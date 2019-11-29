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
		oid := record.SessionId
		if oid == nil {
			fmt.Println("the expected Session ID is nil:[", record.ToCombinedLog(), "].")
			break
		}
		raw, _err := colGeneralRequests.FindOne(newCrudContext(), bson.M{"_id": oid}).DecodeBytes()
		if _err != nil {
			fmt.Println("failed to find the target general request by the previous session ID:", err)
			break
		}
		res, _err := colLandingRequests.InsertOne(newCrudContext(), raw)
		if _err != nil {
			fmt.Println("failed to insert the general request found by the previous session ID:", err, res)
		}
	case secoo.LevelFollowingTimeRequest:
		// ID = NewID(); SID = NextSessionID;
		record.SessionId = &record.Id
		record.Id = dxb.NewObjectId()
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
