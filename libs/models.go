package libs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/zhanbei/static-server/conf"
)

type IRecorder interface {
	NewInstance(start time.Time, realIp string, req *http.Request, code int, header http.Header) IRecord
	DoRecord(record IRecord) error
}

type IRecord interface {
	// Do serialize the record.
	Save() error
	// Print something.
	Log()

	ToCombinedLog() string
}

type Recorder struct {
	*conf.OptionLogger
}

func NewRecorder(ops *conf.OptionLogger) *Recorder {
	return &Recorder{ops}
}

func twoWriters(stdout bool, file *os.File) io.Writer {
	if !stdout {
		if file == nil {
			return nil
		} else {
			return file
		}
	} else {
		if file == nil {
			return os.Stdout
		} else {
			return io.MultiWriter(os.Stdout, file)
		}
	}
}

// FIX-ME The strategy of writing to stdout synchronously and writing to file asynchronously may be applied.
func (m *Recorder) DoRecord(record IRecord) error {
	target := twoWriters(m.Stdout, m.LogWriter)
	_, err := m.Record(target, record)
	return err
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

func (m *Recorder) NewInstance(start time.Time, realIp string, req *http.Request, code int, header http.Header) IRecord {
	return &Record{
		NewDevice(realIp, req.UserAgent()),
		NewRequest(req),
		NewResponse(code, header, time.Since(start)),
		GetCurrentMilliseconds(),
	}
}

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
