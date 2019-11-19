package conf

import (
	"os"
)

type LoggerFormat string

const (
	// Default Log Format = #LoggerFormatExtended
	LoggerFormatText LoggerFormat = "text"

	LoggerFormatCommon LoggerFormat = "common"

	LoggerFormatCombined LoggerFormat = "combined"

	LoggerFormatVirtualHosts LoggerFormat = "vhosts"

	LoggerFormatExtended LoggerFormat = "extended"
	// JSON with Rich Data
	LoggerFormatJson LoggerFormat = "json"
)

var _formats = make(map[LoggerFormat]bool, 0)

func init() {
	_formats[LoggerFormatText] = true
	_formats[LoggerFormatCommon] = true
	_formats[LoggerFormatCombined] = true
	_formats[LoggerFormatVirtualHosts] = true
	_formats[LoggerFormatExtended] = true
	_formats[LoggerFormatJson] = true
}

func IsLoggerFormatValid(m LoggerFormat) bool {
	return _formats[m]
}

type LoggerTarget = string

const (
	LoggerTargetStdout LoggerTarget = "stdout"

	LoggerTargetFile LoggerTarget = "${file}"

	LoggerTargetDir LoggerTarget = "${dir}"
)

type OptionLoggers = []*OptionLogger

// Logger
type OptionLogger struct {
	Enabled bool `json:"enabled"`
	// text | json
	Format LoggerFormat `json:"format"`
	// A log file per host.
	// You are recommended to turn it off, if the hosts received kinds of arbitrary or dynamic,
	// which means that infinite number of logging files may be created.
	// FIX-ME Group similar hosts into a single file may be supported.
	PerHost bool `json:"perHost"`
	// Whether to print the content to be logged to stdout.
	Stdout bool `json:"stdout"`
	// The value should be `stdout|${file}` if the `perHost` option is false.
	// The value should be `${dir}` if the `perHost` option is true.
	Target LoggerTarget `json:"target"`
	// The recorder instance to be used to record requests logs.
	//Recorder IRecorder `json:"-"`

	LogWriter *os.File `json:"-"`
}

func NewLogger(format LoggerFormat, stdout bool, target string) *OptionLogger {
	return &OptionLogger{true, format, false, stdout, target, nil}
}

func (m *OptionLogger) IsValid() bool {
	if doPass(m.Enabled) {
		return true
	}
	if !IsLoggerFormatValid(m.Format) {
		return false
	}
	if strictMode() && !exist(m.Target) && !m.Stdout {
		return false
	}
	return true
}

func (m *OptionLogger) ValidateRequiredResources() error {
	if !m.Enabled {
		return nil
	}
	if !exist(m.Target) {
		return nil
	}
	writer, err := GetFileWriter(m.Target)
	m.LogWriter = writer
	return err
}
