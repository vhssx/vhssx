package configures

import (
	"strings"
)

type SiteConfigure struct {
	DirectoryListing *bool `json:"directoryListing"`

	NoTrailingSlash *bool `json:"noTrailingSlash"`
	// Dynamic routing, from a group of paths to a fixed target.
	RouteMappings *RouteMappings `json:"routeMappings"`

	PrivatePages []string `json:"privatePages"`
	// Allow all robots or disallow.
	//RobotRules []string `json:"robotRules"`
}

type RouteMappings = []RouteMapping
type RouteMapping struct {
	Regexp string `json:"regexp"`

	Target string `json:"target"`
}

func (m *SiteConfigure) IsPrivate(path string) bool {
	if m.PrivatePages == nil || len(m.PrivatePages) == 0 {
		return false
	}
	for _, page := range m.PrivatePages {
		// FIX-ME Use a more efficient strategy to check whitelist, like using map.
		if strings.HasPrefix(path, page) {
			return true
		}
	}
	return false
}
