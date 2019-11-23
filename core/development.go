package core

import (
	"net/http"
	"strings"
)

func TrimSuffixDomainForDevelopment(next http.Handler, suffix string) http.Handler {
	if !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		matched := strings.HasSuffix(r.Host, suffix)
		if matched {
			r.Host = r.Host[:len(r.Host)-len(suffix)]
		}
		next.ServeHTTP(w, r)
	})
}
