package configures

import (
	"strings"
)

const PrefixSpecialSites = "_."

type StaticSites = []*StaticSite

type StaticSite struct {
	// The domain name for the target site.
	Name string `json:"name"`

	Configure *SiteConfigure `json:"configure"`
}

// Is the site special.
func NewSite(name string, conf *SiteConfigure) (site *StaticSite, t bool) {
	if strings.HasPrefix(name, PrefixSpecialSites) {
		name = name[2:]
		t = true
	}
	site = &StaticSite{name, conf}
	return
}
