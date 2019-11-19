package conf_test

import (
	"strconv"
	"testing"

	. "github.com/zhanbei/static-server/conf"
)

func TestConfigure_IsValid(t *testing.T) {
	// Invalid configures example.
	validateConfigures(t, "Validating the false values.", &Configure{
		".",
		"8080",
		&ServerOptions{true, false, true, true},
		[]*OptionLogger{
			NewLogger("json", false, ""),
			&OptionLogger{false, "bson", false, false, "", nil},
		},
		&MongoDbOptions{true, "mongodb://127.0.0.1:27017", "", ""},
		&OptionLoggerGorilla{false, "text", false, "", nil},
	}, false)

	// Valid configures example.
	validateConfigures(t, "Validating the true values.", &Configure{
		".",
		"8080",
		NewDefaultServerOptions(),
		[]*OptionLogger{
			NewLogger("json", false, "logs/whatever.log"),
		},
		&MongoDbOptions{true, "mongodb://127.0.0.1:27017", "vhss", "logging.vhss"},
		&OptionLoggerGorilla{true, "combined", false, "", nil},
	}, true)
}

func validateConfigures(t *testing.T, testCase string, cfg *Configure, expected bool) {
	has := false
	as := func(label string, value bool) {
		if has {
			assert(t, "", label, value, expected)
		} else {
			ok := assert(t, testCase, label, value, expected)
			if !ok {
				has = true
			}
		}
	}
	as("cfg.IsValid()", cfg.IsValid())
	as("cfg.Server.IsValid()", cfg.Server.IsValid())
	for i, logger := range cfg.Loggers {
		as("cfg.Loggers["+strconv.Itoa(i)+"].IsValid()", logger.IsValid())
	}
	as("cfg.MongoDbOptions.IsValid()", cfg.MongoDbOptions.IsValid())
	as("cfg.GorillaOptions.IsValid()", cfg.GorillaOptions.IsValid())

}

func assert(t *testing.T, testCase, label string, value, expected bool) bool {
	if value != expected {
		if testCase != "" {
			t.Log("\n", testCase)
		}
		t.Error("Failure("+label+"):", value, expected)
		return false
	}
	return true
}
