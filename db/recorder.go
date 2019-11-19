package db

import (
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
