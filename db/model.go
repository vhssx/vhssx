package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zhanbei/static-server/libs"
)

type Recorder struct{}

func (m *Recorder) NewInstance(device *libs.Device, request *libs.Request, response *libs.Response) libs.IRecord {
	return &Record{NewObjectId(), device, request, response, libs.GetCurrentMilliseconds()}
}

// The record of a single request, with device and response.
// Naming: Record, RequestRecord, ServiceRecords
type Record struct {
	Id ObjectId `json:"_id"`

	Device *libs.Device `json:"device"`

	Request *libs.Request `json:"req"`

	Response *libs.Response `json:"res"`

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
	fmt.Println(m.ToCombinedLog())
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
