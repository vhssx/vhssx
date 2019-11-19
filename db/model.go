package db

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/recorder"
)

var _ recorder.IRecorder = (*Recorder)(nil)

type Recorder struct {
	ops *conf.MongoDbOptions
}

func NewRecorder(ops *conf.MongoDbOptions) *Recorder {
	return &Recorder{ops}
}

func (m *Recorder) DoRecord(start, end time.Time, realIp string, req *http.Request, code int, header http.Header) error {
	record := &Record{
		NewObjectId(),
		recorder.NewDevice(realIp, req.UserAgent()),
		recorder.NewRequest(req),
		recorder.NewResponse(code, header, end.Sub(start)),
		recorder.GetMilliseconds(start),
	}
	// Asynchronously serialize the record to database( to save time).
	go InsertRecordWithErrorProcessed(record)
	return nil
}

var _ recorder.IRecord = (*Record)(nil)

// The record of a single request, with device and response.
// Naming: Record, RequestRecord, ServiceRecords
type Record struct {
	Id ObjectId `json:"_id"`

	Device *recorder.Device `json:"device"`

	Request *recorder.Request `json:"req"`

	Response *recorder.Response `json:"res"`

	Time int64 `json:"time"`
}

func (m *Record) ToCombinedLog() string {
	req := m.Request
	res := m.Response
	return fmt.Sprintf(
		`%s - - [%s] "%s %s %s" %d %d "%s" "%s"`,
		m.Device.Ip, time.Unix(m.Time/1000, 0).Format("02/Jan/2006:15:04:05 -0700"), req.Method, req.Path, req.Proto, res.Code, res.ContentLength, req.Referer, m.Device.UserAgent,
	)
}
