package conf

import (
	"os"
)

type GorillaFormat string

const (
	GorillaFormatCommonLog GorillaFormat = "common"

	GorillaFormatCombinedLog GorillaFormat = "combined"
)

type OptionLoggerGorilla struct {
	Enabled bool `json:"enabled"`
	// "combined"
	Format GorillaFormat `json:"format"`

	Stdout bool `json:"stdout"`

	Target LoggerTarget `json:"target"`

	LogWriter *os.File `json:"-"`
}

func (m *OptionLoggerGorilla) IsValid() bool {
	if doPass(m.Enabled) {
		return true
	}
	if m.Format != GorillaFormatCommonLog && m.Format != GorillaFormatCombinedLog {
		return false
	}
	if strictMode() && !exist(m.Target) && !m.Stdout {
		return false
	}
	return true
}
