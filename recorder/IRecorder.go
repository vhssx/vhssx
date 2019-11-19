package recorder

import (
	"net/http"
	"time"
)

type IRecorder interface {
	DoRecord(start time.Time, realIp string, req *http.Request, code int, header http.Header) error
}

// It is responsible for recorder to take over everything, logging to stdout or file. :(
// No need for record for now!
type IRecord interface {
	ToCombinedLog() string
}
