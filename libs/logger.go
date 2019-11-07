package libs

import (
	"encoding/json"
	"fmt"
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

func StructuredLoggingHandler(next http.Handler, ops *ServerOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		lrw := NewResponseLoggingWriter(w)
		defer func(start time.Time) {
			// Wondering the differences between the deferred function and direct put below `next.ServeHTTP(w, req)`?
			ip := req.RemoteAddr
			if ops.TrustProxyIp {
				ip = getRemoteIp(req.Header, ip)
			}
			record := NewRecord(
				NewDevice(ip, req.UserAgent()),
				NewRequest(req),
				NewResponse(lrw.StatusCode, w.Header(), time.Since(start)),
			)
			// FIXME Write to database or stdout, following the configures.
			bts, _ := json.Marshal(record)
			fmt.Println(record.ToCombinedLog())
			fmt.Println("-+>", string(bts))
		}(time.Now())
		next.ServeHTTP(lrw, req)
	}
}
