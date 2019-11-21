package sites

import (
	"github.com/zhanbei/static-server/configures"
)

var _Global *configures.ModularSite
var _Modular configures.ModularSites
var _Regular configures.RegularSites
var mModular, mRegular configures.MapRegularSites

func RefreshSites(rootDir string) {
	_Global, _Modular, _Regular = configures.ScanSites(rootDir)
	m := make(configures.MapRegularSites, 0)

	mModular = m
}

func GetCachedRegularSite(host string) *configures.RegularSite {
	for _, site := range _Regular {
		if site.Name == host {
			return site
		}
	}
	return nil
}

func GetCachedModularSite(host string) *configures.ModularSite {
	for _, site := range _Modular {
		if site.IsRootDomain(host) {
			return site
		}
	}
	return nil
}
