package libs

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// The record of a single request, with device and response.
// Naming: Record, RequestRecord, ServiceRecords
type Record struct {
	Device *Device `json:"device"`

	Request *Request `json:"req"`

	Response *Response `json:"res"`

	Time int64 `json:"time"`
}

func NewRecord(device *Device, request *Request, response *Response) *Record {
	return &Record{device, request, response, GetCurrentMilliseconds()}
}

func (m *Record) ToCombinedLog() string {
	req := m.Request
	res := m.Response
	return fmt.Sprintf(
		`%s - - [%s] "%s %s %s" %d %d "%s" "%s"`,
		m.Device.Ip, time.Unix(m.Time/1000, 0).String(), req.Method, req.Path, req.Proto, res.Code, res.ContentLength, req.Referer, m.Device.UserAgent,
	)
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

type Response struct {
	Code int `json:"code,omitempty"`
	//Status string `json:"status,omitempty"`
	ContentLength int64 `json:"length"`

	ContentType string `json:"type,omitempty"`
	//Compressed bool `json:"compressed,omitempty"`
	// Time spent in nanoseconds.
	Duration int64 `json:"duration"`
	// Readable duration.
	Performance string `json:"_duration"`
}

func getContentLength(header http.Header) int64 {
	_length := header.Get("Content-Length")
	if _length == "" {
		return 0
	}
	length, err := strconv.ParseInt(_length, 10, 64)
	if err != nil {
		fmt.Println("Failed to decode the Content-Length:", )
	}
	return length
}

func NewResponse(header http.Header, duration time.Duration) *Response {
	if header != nil {
		return &Response{ContentLength: getContentLength(header), ContentType: header.Get("Content-Type"), Duration: duration.Nanoseconds(), Performance: duration.String()}
	}
	return &Response{Duration: duration.Nanoseconds(), Performance: duration.String()}
	//return &Response{
	//	res.StatusCode, // res.Status,
	//	res.ContentLength, res.Header.Get("Content-Type"), // !res.Uncompressed,
	//	duration.Nanoseconds(), duration.String(),
	//}
}
