package core

import (
	"net/http"

	"github.com/zhanbei/serve-static"
	"github.com/zhanbei/static-server/sites"
)

// Another pattern is to create server for all existed sites, with standalone configurations.
func VirtualHostStaticHandler(ss *servestatic.FileServer) http.Handler {
	sites.RefreshSites(ss.RootDir)
	return &mStaticServer{ss}
}

type mStaticServer struct {
	ss *servestatic.FileServer
}

// Fallthrough Resources:
//
// 1. Cached Regular Site --> 2. Real Universal Site --> 3. Cached Chained Modular Site --> 4. Not Found 404
//
// Fallthrough 404/etc:
//
// 1. Cached Regular Site 404 --> 2. Cached Chained Modular Site --> 3. Not Found 404
func (m *mStaticServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch the cached target site configures;
	// The regular site will be preferred, then is the modular site.
	site := sites.GetCachedRegularSite(r.Host)
	if site != nil {
		site.ServeHTTP(w, r)
		return
	}

	// 2. Serve static files in the normal way.
	m.ss.ServeFiles(w, r, func(location string) {
		// 3. Get the cached modular sites.
		module := sites.GetCachedModularSite(r.Host)
		if module != nil {
			module.ServeHTTP(w, r)
			return
		}
		// 4. Real Not Found
		w.WriteHeader(http.StatusNotFound)
	})
}
