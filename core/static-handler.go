package core

import (
	"fmt"
	"net/http"

	"github.com/zhanbei/serve-static"
	"github.com/zhanbei/static-server/configures"
)

var _Global *configures.StaticSite
var _Modular, _Sites configures.StaticSites
var mModular, mSite configures.MapStaticSites

func RefreshSites(rootDir string) {
	_Global, _Modular, _Sites = configures.ScanSites(rootDir)
	m := make(configures.MapStaticSites, 0)

	mModular = m
}

func GetModularSite(host string) *configures.StaticSite {
	for _, site := range mModular {
		if site.IsRootDomain(host) {
			return site
		}
	}
	return nil
}

// Another pattern is to create server for all existed sites, with standalone configurations.
func VirtualHostStaticHandler(ss *servestatic.FileServer) http.Handler {
	RefreshSites(ss.RootDir)
	return &mStaticServer{ss}
}

type mStaticServer struct {
	ss *servestatic.FileServer
}

// - Take care of hosts in the development mode;
// - Find the target site configures;
// - Serve the static file of target site;
// - Path mapping rendering
// - Fallthrough
func (m *mStaticServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	site := mSite[r.Host]

	// 1. Filters for private pages to protect whitelist(hidden resources).
	if site != nil && site.Configure != nil && site.Configure.IsPrivate(r.URL.Path) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 2. Serve static files in the normal way.
	m.ss.ServeFiles(w, r, func(location string) {
		// 3. Use routes mappers.

		// Find the target resource failed, hence fallthrough for custom pages.
		// a. Site/target --> b. Scope/target --> c. Global/target
		// a. Site/404.html --> b. Scope/404.html --> c. Global/404.html
		module := GetModularSite(r.Host)
		if module == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		module.StaticServer.ServeFiles(w, r, func(resolvedLocation string) {
			w.WriteHeader(http.StatusNotFound)
			// Write the custom 404 page.
			fmt.Println("Requested file is not found:", "http://"+r.Host+r.RequestURI, "Resolved File Location:", location)
		})
	})
}
