package db

import (
	"fmt"
	"time"

	"github.com/zhanbei/dxb"
	"github.com/zhanbei/static-server/recorder"
	"github.com/zhanbei/static-server/secoo"
)

var _ recorder.IRecord = (*Record)(nil)

// The record of a single request, with device and response.
// Naming: Record, RequestRecord, ServiceRecords
type Record struct {
	Id dxb.ObjectId `json:"_id"`

	SessionId *dxb.ObjectId `json:"_sid,omitempty"`

	Session *secoo.SessionCookieStore `json:"session,omitempty"`

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
