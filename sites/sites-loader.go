package sites

import (
	"strings"

	"github.com/zhanbei/static-server/configures"
)

var _Global *configures.ModularSite
var _Modular configures.ModularSites
var _Regular configures.RegularSites

var mModular configures.MapModularSites
var mRegular configures.MapRegularSites

func RefreshSites(rootDir string) {
	_Global, _Modular, _Regular = configures.ScanSites(rootDir)

	sites := make(configures.MapRegularSites, 0)
	for _, site := range _Regular {
		sites[site.Name] = site
	}
	mRegular = sites

	modular := make(configures.MapModularSites, 0)
	for _, module := range _Modular {
		modular[module.Name] = module
	}
	mModular = modular
}

func GetCachedRegularSite(host string) *configures.RegularSite {
	return mRegular[host]
}

func GetCachedModularSite(host string) *configures.ModularSite {
	if mModular[host] != nil {
		// Optimization for bare site.
		return mModular[host]
	}
	// Get a trial for optimization -- the root of all evil.
	hosts := strings.SplitN(host, ".", 2)
	if len(hosts) > 1 && mModular[hosts[1]] != nil {
		// Optimization for simple sub domains.
		return mModular[hosts[1]]
	}
	// O(n) is acceptable. :)
	for _, site := range _Modular {
		//if site.IsRootDomain(host) {
		if strings.HasSuffix(host, "."+site.Name) {
			return site
		}
	}
	return _Global
}
