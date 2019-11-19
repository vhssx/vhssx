package recorder

import (
	"net/http"
	"time"
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
