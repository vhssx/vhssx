package recorder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var _ IRecord = (*Record)(nil)

// The record of a single request, with device and response.
// Naming: Record, RequestRecord, ServiceRecords
type Record struct {
	Device *Device `json:"device"`

	Request *Request `json:"req"`

	Response *Response `json:"res"`

	Time int64 `json:"time"`
}

func (m *Record) ToCombinedLog() string {
	req := m.Request
	res := m.Response
	return fmt.Sprintf(
		`%s - - [%s] "%s %s %s" %d %d "%s" "%s" "%s"`,
		m.Device.Ip, time.Unix(m.Time/1000, 0).Format("02/Jan/2006:15:04:05 -0700"), req.Method, req.Path, req.Proto, res.Code, res.ContentLength, req.Referer, m.Device.UserAgent,
		time.Duration(m.Response.Duration).String(),
	)
}

func (m *Record) Log() {
	fmt.Println(m.ToCombinedLog(), time.Duration(m.Response.Duration).String())
	bts, _ := json.Marshal(m)
	fmt.Println("-+>", string(bts))
}

func (m *Record) Save() error {
	return nil
}

type Device struct {
	// Trust proxy by default.
	Ip string `json:"ip"`

	UserAgent string `json:"info"`
}

func NewDevice(ip, userAgent string) *Device {
	return &Device{ip, userAgent}
}

// See combined request logging.
type Request struct {
	// Omit the value if the value is "HTTP/1.1".
	Proto string `json:"proto,omitempty"`

	Method string `json:"method"`

	Host string `json:"host"`

	Path string `json:"path"`

	Referer string `json:"referer,omitempty"`

	ContentLength int64 `json:"length"`
}

func NewRequest(req *http.Request) *Request {
	//un, _, _ := req.BasicAuth()
	return &Request{
		req.Proto,
		req.Method, req.Host, req.URL.Path,
		req.Referer(),
		req.ContentLength,
	}
}

// The summaries of the response.
type Response struct {
	Code int `json:"code,omitempty"`
	//Status string `json:"status,omitempty"`
	ContentLength int64 `json:"length"`

	ContentType string `json:"type,omitempty"`
	//Compressed bool `json:"compressed,omitempty"`
	// Time spent in nanoseconds.
	Duration int64 `json:"duration"`
	// Readable duration.
	//Performance string `json:"_duration"`
}

func getContentLength(header http.Header) int64 {
	if header == nil {
		return 0
	}
	_length := header.Get("Content-Length")
	if _length == "" {
		return 0
	}
	length, err := strconv.ParseInt(_length, 10, 64)
	if err != nil {
		fmt.Println("Failed to decode the Content-Length:")
	}
	return length
}

func getContentType(header http.Header) string {
	if header == nil {
		return ""
	}
	return header.Get("Content-Type")
}

func NewResponse(code int, header http.Header, duration time.Duration) *Response {
	return &Response{code, getContentLength(header), getContentType(header), duration.Nanoseconds()}
}
