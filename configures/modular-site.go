package configures

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zhanbei/serve-static"
)

const PrefixSpecialSites = "_."

type ModularSites = []*ModularSite
type MapModularSites = map[string]*ModularSite

// A modular site is the root of a group regular sites, for resources to fall through.
type ModularSite struct {
	// The domain name for the target site.
	Name string `json:"name"`

	Configure *SiteConfigure `json:"configure"`

	StaticServer *servestatic.FileServer
}

// Is the site special.
func NewModularSite(name string, conf *SiteConfigure) *ModularSite {
	return &ModularSite{name, conf, nil}
}

func (m *ModularSite) IsRootDomain(host string) bool {
	if !strings.HasSuffix(host, m.Name) {
		return false
	}
	return strings.HasSuffix(host, "."+m.Name) || host == m.Name
}

func (m *ModularSite) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Filters for private pages to protect whitelist(hidden resources).
	if m.Configure != nil && m.Configure.IsPrivate(r.URL.Path) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	m.StaticServer.ServeFiles(w, r, func(resolvedLocation string) {
		w.WriteHeader(http.StatusNotFound)
		// Write the custom 404 page.
		fmt.Println("Requested file is not found:", "http://"+r.Host+r.RequestURI, "Resolved File Location:", resolvedLocation)
	})
}
