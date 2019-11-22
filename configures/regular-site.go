package configures

import (
	"net/http"

	"github.com/zhanbei/serve-static"
	"github.com/zhanbei/static-server/helpers/terminator"
)

type RegularSites = []*RegularSite
type MapRegularSites = map[string]*RegularSite

type RegularSite struct {
	// The domain name for the target site.
	Name string `json:"name"`

	Configure *SiteConfigure `json:"configure"`
	// [CACHE] The parent modular site to fallthrough.
	ModularSite *ModularSite `json:"-"`
	// A  static server to serve resources in the no-trailing-slash mode.
	StaticServer *servestatic.FileServer
}

func NewRegularSite(name, dirSiteRoot string, conf *SiteConfigure) *RegularSite {
	// FIX-ME Setup server following configures.
	server, err := servestatic.NewFileServer(dirSiteRoot, false)
	if err != nil {
		terminator.ExitWithPreLaunchServerError(err, "Setting up the static server for site ["+name+"]("+dirSiteRoot+") failed!")
	}
	return &RegularSite{name, conf, nil, server}
}

// - routes mappers.
func (m *RegularSite) ServeHTTP(w http.ResponseWriter, r *http.Request, notFound func()) {
	// 1. Filters for private pages to protect whitelist(hidden resources).
	if ServeDynamicContents(w, r, m.Configure, m.StaticServer, notFound) {
		return
	}
	// Falling through target resources:
	// 1. Cached Regular Site --> 2. Real Universal Site --> 3. Cached Chained Modular Site --> 4. Not Found 404
	m.StaticServer.ServeFiles(w, r, func(resolvedLocation string) {
		// Not found the target resource.
		if m.ModularSite == nil {
			notFound()
			return
		}
		m.ModularSite.ServeHTTP(w, r, func() {
			// Not found the target resource by parent modular sites.
			notFound()
		})
	})
}

// Responding 404:
// 1. Cached Regular Site 404 --> 2. Cached Chained Modular Site --> 3. Not Found 404
func (m *RegularSite) RespondingCustom404(w http.ResponseWriter, r *http.Request, prev func()) {
	// 1. Cached Regular Site 404
	handled, _ := Serve404Page(w, m.StaticServer, Page404Path)
	// FIX-ME If error is encountered, just call prev() and return.
	if handled {
		return
	}

	if m.ModularSite != nil {
		// 2. Cached Chained Modular Site
		m.ModularSite.Responding404(w, r, prev)
		return
	}

	// 3. Not Found 404
	prev()
}
