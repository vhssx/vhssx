package db

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zhanbei/static-server/recorder"
)

type Recorder struct{}

func NewRecorder() *Recorder {
	return new(Recorder)
}

func (m *Recorder) NewInstance(start time.Time, realIp string, req *http.Request, code int, header http.Header) recorder.IRecord {
	return &Record{
		NewObjectId(),
		recorder.NewDevice(realIp, req.UserAgent()),
		recorder.NewRequest(req),
		recorder.NewResponse(code, header, time.Since(start)),
		recorder.GetCurrentMilliseconds(),
	}
}

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

func (m *Record) Log() {
	fmt.Println(m.ToCombinedLog(), time.Duration(m.Response.Duration).String())
	bts, _ := json.Marshal(m)
	fmt.Println("-+>", string(bts))
}

func (m *Record) Insert() error {
	_, err := col.InsertOne(newCrudContext(), m)
	return err
}

func (m *Record) Save() error {
	go m.Insert()
	return nil
}
