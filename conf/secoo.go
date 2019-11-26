package conf

type OptionsSessionCookie struct {
	Enabled bool `json:"enabled"`

	Secret string `json:"secret"`
	// Whether to share cookies across all sub domains.
	AllSubDomains bool `json:"withSubDomains"`
	// Which strategy to use.
	Strategy string `json:"strategy"`
}

func (m *OptionsSessionCookie) IsValid() bool {
	if m.Enabled {
		return exist(m.Secret)
	}
	return true
}
