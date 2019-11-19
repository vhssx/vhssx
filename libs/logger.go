package libs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zhanbei/static-server/conf"
)

func getRemoteIp(header http.Header, original string) string {
	ip := header.Get("x-remote-addr")
	if ip != "" {
		return ip
	}
	ip = header.Get("x-forwarded-for")
	if ip != "" {
		return ip
	}
	ip = header.Get("ip")
	if ip != "" {
		return ip
	}
	return original
}

//@see https://ndersson.me/post/capturing_status_code_in_net_http/
type ResponseLoggingWriter struct {
	http.ResponseWriter

	StatusCode int
}

func NewResponseLoggingWriter(w http.ResponseWriter) *ResponseLoggingWriter {
	return &ResponseLoggingWriter{w, http.StatusOK}
}

func (m *ResponseLoggingWriter) WriteHeader(code int) {
	m.StatusCode = code
	m.ResponseWriter.WriteHeader(code)
}

func StructuredLoggingHandler(next http.Handler, cfg *conf.Configure) http.HandlerFunc {
	ops := cfg.Server

	recorders := make([]IRecorder, 0)

	for _, logger := range cfg.Loggers {
		if !logger.Enabled {
			continue
		}
		if logger.PerHost {
			fmt.Println("Not supported logger#PerHost:", logger)
			continue
		}
		recorder := NewRecorder(logger)
		recorders = append(recorders, recorder)
	}

	mon := cfg.MongoDbOptions
	gor := cfg.GorillaOptions
	if (gor == nil || !gor.Enabled) && (mon == nil || !mon.Enabled) && len(recorders) == 0 { // <= 0 {
		// Add a default console(stdout) logger when there is no logger configured!
		logger := conf.NewLogger(conf.LoggerFormatText, true, "")
		recorder := NewRecorder(logger)
		recorders = append(recorders, recorder)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		lrw := NewResponseLoggingWriter(w)
		defer func(start time.Time) {
			// Wondering the differences between the deferred function and direct put below `next.ServeHTTP(w, req)`?
			ip := req.RemoteAddr
			if ops.TrustProxyIp {
				ip = getRemoteIp(req.Header, ip)
			}

			for _, recorder := range recorders {
				record := recorder.NewInstance(start, ip, req, lrw.StatusCode, w.Header())
				_ = recorder.DoRecord(record)
			}
		}(time.Now())
		next.ServeHTTP(lrw, req)
	}
}
