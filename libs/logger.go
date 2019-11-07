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

func StructuredLoggingHandler(next http.Handler, ops *ServerOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func(start time.Time) {
			ip := req.RemoteAddr
			if ops.TrustProxyIp {
				ip = getRemoteIp(req.Header, ip)
			}
			record := NewRecord(
				NewDevice(ip, req.UserAgent()),
				NewRequest(req),
				NewResponse(w.Header(), time.Since(start)),
			)
			// FIXME Write to database or stdout, following the configures.
			bts, _ := json.Marshal(record)
			fmt.Println(record.ToCombinedLog())
			fmt.Println("-+>", string(bts))
		}(time.Now())
		next.ServeHTTP(w, req)
	}
}
