package recorder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zhanbei/static-server/conf"
)

type Recorder struct {
	*conf.OptionLogger
}

func NewRecorder(ops *conf.OptionLogger) *Recorder {
	return &Recorder{ops}
}

func (m *Recorder) Record(target io.Writer, record IRecord) (int, error) {
	if target == nil {
		return -1, nil
	}
	if m.Format == conf.LoggerFormatText {
		return fmt.Fprintln(target, record.ToCombinedLog())
	} else if m.Format == conf.LoggerFormatJson {
		bts, err := json.Marshal(record)
		if err != nil {
			return -1, err
		}
		bts = append(bts, '\n')
		return target.Write(bts)
	} else {
		return -1, errors.New("unsupported logger format: " + string(m.Format))
	}
}

// FIX-ME The strategy of writing to stdout synchronously and writing to file asynchronously may be applied.
func (m *Recorder) DoRecord(start time.Time, realIp string, req *http.Request, code int, header http.Header) error {
	record := &Record{
		NewDevice(realIp, req.UserAgent()),
		NewRequest(req),
		NewResponse(code, header, time.Since(start)),
		GetCurrentMilliseconds(),
	}
	target := twoWriters(m.Stdout, m.LogWriter)
	_, err := m.Record(target, record)
	return err
}
