package configures

import (
	"net/http"
	"strings"

	"github.com/zhanbei/serve-static"
	"github.com/zhanbei/static-server/helpers/terminator"
)

const PrefixSpecialSites = "_."

type ModularSites = []*ModularSite
type MapModularSites = map[string]*ModularSite

// A modular site is the root of a group regular sites, for resources to fall through.
type ModularSite struct {
	// The domain name for the target site.
	Name string `json:"name"`

	Configure *SiteConfigure `json:"configure"`
	// [CACHE] The parent modular site to fallthrough.
	ModularSite *ModularSite `json:"-"`

	StaticServer *servestatic.FileServer
}

// Is the site special.
func NewModularSite(name, dirSiteRoot string, conf *SiteConfigure) *ModularSite {
	// FIX-ME Setup server following configures.
	server, err := servestatic.NewFileServer(dirSiteRoot, false)
	if err != nil {
		terminator.ExitWithPreLaunchServerError(err, "Setting up the static server for site ["+name+"]("+dirSiteRoot+") failed!")
	}
	return &ModularSite{name, conf, nil, server}
}

func (m *ModularSite) IsRootDomain(host string) bool {
	if !strings.HasSuffix(host, m.Name) {
		return false
	}
	return strings.HasSuffix(host, "."+m.Name) || host == m.Name
}

func (m *ModularSite) ServeHTTP(w http.ResponseWriter, r *http.Request, notFound func()) {
	// 1. Filters for private pages to protect whitelist(hidden resources).
	if ServeDynamicContents(w, r, m.Configure, m.StaticServer, notFound) {
		return
	}

	m.StaticServer.ServeFiles(w, r, func(resolvedLocation string) {
		// Not found the target resource.
		if m.ModularSite == nil {
			notFound()
			return
		}
		m.ModularSite.ServeHTTP(w, r, notFound)
	})
}

// Responding 404 through chained modular sites.
func (m *ModularSite) Responding404(w http.ResponseWriter, r *http.Request, prev func()) {
	handled, _ := Serve404Page(w, m.StaticServer, Page404Path)
	// FIX-ME If error is encountered, just call prev() and return.
	if handled {
		return
	}

	if m.ModularSite != nil {
		m.ModularSite.Responding404(w, r, prev)
		return
	}

	// Not got handled
	prev()
}
