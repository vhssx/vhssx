package configures

import (
	"strings"

	"github.com/zhanbei/serve-static"
)

const PrefixSpecialSites = "_."

type StaticSites = []*StaticSite
type MapStaticSites = map[string]*StaticSite

type StaticSite struct {
	// The domain name for the target site.
	Name string `json:"name"`

	Configure *SiteConfigure `json:"configure"`

	StaticServer *servestatic.FileServer
}

// Is the site special.
func NewSite(name string, conf *SiteConfigure) (site *StaticSite, t bool) {
	if strings.HasPrefix(name, PrefixSpecialSites) {
		name = name[2:]
		t = true
	}
	site = &StaticSite{name, conf, nil}
	return
}

func (m *StaticSite) IsRootDomain(host string) bool {
	if !strings.HasSuffix(host, m.Name) {
		return false
	}
	return strings.HasSuffix(host, "."+m.Name) || host == m.Name
}
