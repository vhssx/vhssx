package recorder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/helpers/writersHelper"
	"github.com/zhanbei/static-server/utils"
)

var _ IRecorder = (*Recorder)(nil)

type Recorder struct {
	*conf.OptionLogger
}

func NewRecorder(ops *conf.OptionLogger) *Recorder {
	return &Recorder{ops}
}

func (m *Recorder) Record(target io.Writer, record *Record) (int, error) {
	if target == nil {
		return -1, nil
	}
	if m.Format == conf.LoggerFormatJson {
		bts, err := json.Marshal(record)
		if err != nil {
			return -1, err
		}
		n, err := target.Write(append(bts, '\n'))
		if err != nil {
			PrintFailedRecordText(string(bts))
		}
		return n, err
	}

	content := ""
	switch m.Format {
	case conf.LoggerFormatText:
		// Default as extended format.
		content = record.ToExtendedLog()
	case conf.LoggerFormatCommon:
		content = record.ToCommonLog()
	case conf.LoggerFormatCombined:
		content = record.ToCombinedLog()
	case conf.LoggerFormatVirtualHosts:
		content = record.ToVirtualHostsLog()
	case conf.LoggerFormatExtended:
		content = record.ToExtendedLog()
	default:
		// By default, this is not possible.
		content := record.ToExtendedLog()
		PrintFailedRecordText(content)
		return -1, errors.New("unsupported logger format: " + string(m.Format))
	}

	n, err := fmt.Fprintln(target, content)
	if err != nil {
		PrintFailedRecordText(content)
	}
	return n, err
}

// FIX-ME The strategy of writing to stdout synchronously and writing to file asynchronously may be applied.
func (m *Recorder) DoRecord(start, end time.Time, realIp string, req *http.Request, code int, header http.Header) error {
	record := &Record{
		NewDevice(realIp, req.UserAgent()),
		NewRequest(req),
		NewResponse(code, header, end.Sub(start)),
		utils.GetMilliseconds(start),
	}
	target := writersHelper.StdoutVsFileWriter(m.Stdout, m.LogWriter)
	_, err := m.Record(target, record)
	return err
}
