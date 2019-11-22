package configures

import (
	"errors"
	"regexp"
	"strings"

	"github.com/zhanbei/static-server/utils"
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

type RouteMappings = []*RouteMapping
type RouteMapping struct {
	Regexp string `json:"pattern"`

	Target string `json:"target"`

	mRegexp *regexp.Regexp `json:"-"`
}

func (m *SiteConfigure) ValidateRequiredResources() error {
	if m.RouteMappings != nil {
		for _, mapping := range *m.RouteMappings {
			if !utils.NotEmpty(mapping.Target) || !utils.NotEmpty(mapping.Regexp) {
				return errors.New("the target and regexp pattern of the router mapping should not be empty")
			}
			exp, err := regexp.Compile(mapping.Regexp)
			if err != nil {
				return err
			}
			mapping.mRegexp = exp
		}
	}
	return nil
}

func (m *SiteConfigure) GetPotentialMappedTarget(path string) string {
	if m.RouteMappings == nil {
		return ""
	}
	for _, mapping := range *m.RouteMappings {
		if mapping.mRegexp == nil {
			continue
		}
		if mapping.mRegexp.MatchString(path) {
			return mapping.Target
		}
		// Check the path, for remapping.
	}
	return ""
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
