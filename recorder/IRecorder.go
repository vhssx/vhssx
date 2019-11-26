package recorder

import (
	"net/http"
	"time"

	"github.com/zhanbei/static-server/secoo"
)

type IRecorders = []IRecorder

type IRecorder interface {
	DoRecord(start, end time.Time, realIp string, req *http.Request, extra *secoo.SessionCookieStore, code int, header http.Header) error
}

// It is responsible for recorder to take over everything, logging to stdout or file. :(
// No need for record for now!
type IRecord interface {
	ToCombinedLog() string
}
