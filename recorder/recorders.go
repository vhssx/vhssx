package recorder

import (
	"fmt"

	"github.com/zhanbei/static-server/conf"
)

func GetActiveRecorders(loggers conf.OptionLoggers) IRecorders {
	recs := make(IRecorders, 0)

	for _, logger := range loggers {
		if !logger.Enabled {
			continue
		}
		if logger.PerHost {
			fmt.Println("Not supported logger#PerHost:", logger)
			continue
		}
		recs = append(recs, NewRecorder(logger))
	}
	return recs
}

// Get a default recorder if there are no records at all.
func GetDefaultRecorder() IRecorder {
	logger := conf.NewLogger(conf.LoggerFormatText, true, "")
	return NewRecorder(logger)
}

func GetDefaultRecorders() IRecorders {
	return IRecorders{GetDefaultRecorder()}
}
