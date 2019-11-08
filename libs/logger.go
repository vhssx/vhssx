package libs

import (
	"net/http"
	"time"
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

func StructuredLoggingHandler(next http.Handler, ops *ServerOptions, recorder IRecorder) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		lrw := NewResponseLoggingWriter(w)
		defer func(start time.Time) {
			// Wondering the differences between the deferred function and direct put below `next.ServeHTTP(w, req)`?
			ip := req.RemoteAddr
			if ops.TrustProxyIp {
				ip = getRemoteIp(req.Header, ip)
			}
			record := recorder.NewInstance(start, ip, req, lrw.StatusCode, w.Header())
			_ = record.Save()
			// FIXME Write to database or stdout, following the configures.
			record.Log()
		}(time.Now())
		next.ServeHTTP(lrw, req)
	}
}
