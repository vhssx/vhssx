package db

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/recorder"
	"github.com/zhanbei/static-server/secoo"
	"github.com/zhanbei/static-server/utils"
)

var _ recorder.IRecorder = (*Recorder)(nil)

type Recorder struct {
	ops *conf.MongoDbOptions
}

func NewRecorder(ops *conf.MongoDbOptions) *Recorder {
	return &Recorder{ops}
}

func (m *Recorder) DoRecord(start, end time.Time, realIp string, req *http.Request, session *secoo.SessionCookieStore, code int, header http.Header) error {
	var sid ObjectId
	if session != nil {
		_sid, err := NewObjectIdFromHex(session.Id)
		if err != nil {
			fmt.Println("failed to parse the session ID from hex:[", session.Id, "].")
		} else {
			sid = _sid
		}
	}
	record := &Record{
		NewObjectId(),
		// Leave a session ID here, is sufficiently fine.
		sid,
		// Record the session store in the collection of first-time requests.
		session,
		recorder.NewDevice(realIp, req.UserAgent()),
		recorder.NewRequest(req),
		recorder.NewResponse(code, header, end.Sub(start)),
		utils.GetMilliseconds(start),
	}
	// Asynchronously serialize the record to database( to save time).
	go InsertRecordWithErrorProcessed(record)
	return nil
}
