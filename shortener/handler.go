package shortener

import (
	"net/http"
)

var checker CacheDomainChecker
var mapper CacheDomainMapper

// FIX-ME Whether to set session cookie for redirection?
// Yes for now, because they may also be counted and analyzed.
func HandlerUrlDirections(next http.Handler) http.HandlerFunc {

	// Load resources required.
	checker, mapper = DoLoadRedirectionRecords()

	return func(w http.ResponseWriter, req *http.Request) {
		if checker[req.Host] {
			domain := mapper[req.Host]
			handled := domain.ServeHttp(w, req)
			if handled {
				// Got handled.
				return
			}
		}
		next.ServeHTTP(w, req)
	}
}
