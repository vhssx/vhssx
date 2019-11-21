package configures

import (
	"net/http"

	"github.com/zhanbei/serve-static"
)

type RegularSites = []*RegularSite
type MapRegularSites = map[string]*RegularSite

type RegularSite struct {
	// The domain name for the target site.
	Name string `json:"name"`

	Configure *SiteConfigure `json:"configure"`
	// A  static server to serve resources in the no-trailing-slash mode.
	StaticServer *servestatic.FileServer
}

func NewRegularSite(name string, conf *SiteConfigure) *RegularSite {
	return &RegularSite{name, conf, nil}
}

// - routes mappers.
func (m *RegularSite) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Filters for private pages to protect whitelist(hidden resources).
	if m.Configure != nil && m.Configure.IsPrivate(r.URL.Path) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, _ = w.Write([]byte("Hola, World!"))
}
