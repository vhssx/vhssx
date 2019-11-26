package conf

import (
	"regexp"
)

const DefaultCrawlerFilter = "(?i)google|baidu|bot|crawler|spider|crawling"

type OptionsSessionCookie struct {
	Enabled bool `json:"enabled"`
	// The field to receive configuration values.
	Secret string `json:"secret,omitempty"`
	// The real secret that is hidden for serialization.
	RealSecret string `json:"-" bson:"-"`
	// Whether to share cookies across all sub domains.
	AllSubDomains bool `json:"withSubDomains"`
	// Filter off the common crawlers with regexp by user-agent.
	// @see https://support.google.com/webmasters/answer/1061943?hl=en
	Filter string `json:"filter"`
	// [CACHE]
	RegexpFilter *regexp.Regexp `json:"-"`
	// Which strategy to use.
	Strategy string `json:"strategy"`
}

func (m *OptionsSessionCookie) IsValid() bool {
	if !m.Enabled {
		return true
	}
	if m.Filter == " " {
		m.Filter = DefaultCrawlerFilter
	}
	if exist(m.Filter) {
		ex, err := regexp.Compile(m.Filter)
		if err != nil {
			return false
		}
		m.RegexpFilter = ex
	}
	m.RealSecret = m.Secret
	// Keep the real secret same and sound.
	m.Secret = ""
	return exist(m.RealSecret)
}
