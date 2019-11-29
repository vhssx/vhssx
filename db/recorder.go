package db

import (
	"net/http"
	"time"

	"github.com/zhanbei/dxb"
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
	// Session is nil, when the secoo module is disabled.
	var sid dxb.ObjectId
	var xid *dxb.ObjectId
	if session == nil {
		sid = dxb.NewObjectId() // Only generate new Object IDs for the no secoo mode.
	} else {
		sid, xid = session.GetSessionIds()
	}
	// Current the #Id is the #SessionId while the #SessionId is the potential #PreviousSessionId.
	record := &Record{
		sid, // NewObjectId(), There are (crawlers/1st/2nd)times when the ID need no generated again.
		// Leave a session ID here, is sufficiently fine.
		xid, // Keep the session ID and extra ID stored here temporarily.
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
