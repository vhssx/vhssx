package libs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/recorder"
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

	recs := make([]recorder.IRecorder, 0)

	for _, logger := range cfg.Loggers {
		if !logger.Enabled {
			continue
		}
		if logger.PerHost {
			fmt.Println("Not supported logger#PerHost:", logger)
			continue
		}
		recs = append(recs, recorder.NewRecorder(logger))
	}

	mon := cfg.MongoDbOptions
	gor := cfg.GorillaOptions
	if (gor == nil || !gor.Enabled) && (mon == nil || !mon.Enabled) && len(recs) == 0 { // <= 0 {
		// Add a default console(stdout) logger when there is no logger configured!
		logger := conf.NewLogger(conf.LoggerFormatText, true, "")
		recs = append(recs, recorder.NewRecorder(logger))
	}

	return func(w http.ResponseWriter, req *http.Request) {
		lrw := NewResponseLoggingWriter(w)
		defer func(start time.Time) {
			// Wondering the differences between the deferred function and direct put below `next.ServeHTTP(w, req)`?
			ip := req.RemoteAddr
			if ops.TrustProxyIp {
				ip = getRemoteIp(req.Header, ip)
			}

			for _, rec := range recs {
				_ = rec.DoRecord(start, ip, req, lrw.StatusCode, w.Header())
			}
		}(time.Now())
		next.ServeHTTP(lrw, req)
	}
}
