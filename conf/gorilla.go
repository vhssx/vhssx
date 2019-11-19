package conf

import (
	"os"
)

type OptionLoggerGorilla struct {
	Enabled bool `json:"enabled"`
	// "combined"
	Format string `json:"format"`

	Stdout bool `json:"stdout"`

	Target LoggerTarget `json:"target"`

	LogWriter *os.File `json:"-"`
}
