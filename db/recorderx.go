package db

import (
	"fmt"

	"github.com/zhanbei/dxb"
	"github.com/zhanbei/static-server/recorder"
	"github.com/zhanbei/static-server/secoo"
	"go.mongodb.org/mongo-driver/bson"
)

// Calling this func asynchronous is recommended.
func InsertRecordWithErrorProcessed(record *Record) {
	err := InsertRecord(record)
	if err != nil {
		recorder.PrintFailedRecord(record)
	}
}

func InsertRecord(record *Record) (err error) {
	if record.Session == nil {
		// Insert record in the no-session-cookie mode.
		_, err = colGeneralRequests.InsertOne(newCrudContext(), record)
		return
	}

	// Current the #Id is the #SessionId while the #SessionId is the potential #PreviousSessionId.
	switch record.Session.Level {
	case secoo.LevelCrawlerRequest:
		// ID = NewID(); SID = nil;
		record.Session = nil
		_, err = colCrawlerRequests.InsertOne(newCrudContext(), record)
	case secoo.LevelFirstTimeRequest:
		// ID = PreviousSessionID; SID = nil;
		_, err = colGeneralRequests.InsertOne(newCrudContext(), record)
	case secoo.LevelSecondTimeRequest:
		// ID = NextSessionID; SID = PreviousSessionID;
		record.Session = nil
		_, err = colValidatingRequests.InsertOne(newCrudContext(), record)
		// Pull the corresponding requests from first level requests out.
		PullOverLandingRequests(record, record.SessionId)
	case secoo.LevelFollowingTimeRequest:
		// ID = NewID(); SID = NextSessionID;
		fmt.Println("record.Id, record.SessionId:", record.Id, record.SessionId)
		record.SessionId = &record.Id
		record.Id = dxb.NewObjectId()
		fmt.Println("record.Id, record.SessionId:", record.Id, record.SessionId)
		record.Session = nil
		_, err = colValidatedRequests.InsertOne(newCrudContext(), record)
	default:
		// ID = NewID(); SID = nil;
		record.SessionId = &record.Id
		// Create a new ID for whatever safety reason.
		record.Id = dxb.NewObjectId()
		_, err = colUnknownRequests.InsertOne(newCrudContext(), record)
	}
	return err
}

func PullOverLandingRequests(record *Record, initialSessionId *dxb.ObjectId) {
	if initialSessionId == nil {
		fmt.Println("the expected Session ID is nil:[", record.ToCombinedLog(), "].")
		return
	}
	raw, err := colGeneralRequests.FindOne(newCrudContext(), bson.M{"_id": initialSessionId}).DecodeBytes()
	if err != nil {
		fmt.Println("failed to find the target general request by the previous session ID:", err)
		return
	}
	res, err := colLandingRequests.InsertOne(newCrudContext(), raw)
	if err != nil {
		fmt.Println("failed to insert the general request found by the previous session ID:", err, res)
	}
}
