package conf

import "github.com/zhanbei/static-server/utils"

type AppMode = string

const (
	ModeDevelopment AppMode = "development"

	ModeProduction AppMode = "production"
)

const DefaultDevDomainSuffix = "loxal.me"

type OptionsApp struct {
	Mode string `json:"mode"`
	// The domain to be used for the development of virtual hosts, if in the development & virtual hosting modes.
	// By default, it is loxal.me
	// Support multiple, if needed.
	DevDomainSuffix string `json:"devDomain"`
}

func NewDefaultAppOptions() *OptionsApp {
	return &OptionsApp{ModeProduction, DefaultDevDomainSuffix}
}

func (m *OptionsApp) IsValid() bool {
	if !m.IsInDevelopmentMode() && m.Mode != ModeProduction {
		return false
	}
	if m.Mode == ModeDevelopment {
		if !utils.NotEmpty(m.DevDomainSuffix) {
			m.DevDomainSuffix = DefaultDevDomainSuffix
		}
	}
	return true
}

func (m *OptionsApp) IsInDevelopmentMode() bool {
	return m.Mode == ModeDevelopment
}
