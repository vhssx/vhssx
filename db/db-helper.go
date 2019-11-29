package db

import (
	"fmt"

	"github.com/zhanbei/dxb"
	"github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/recorder"
	"github.com/zhanbei/static-server/secoo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	app *conf.MongoDbOptions

	mDbClient *mongo.Client

	colCrawlerRequests *mongo.Collection

	colGeneralRequests *mongo.Collection

	colLandingRequests *mongo.Collection

	colValidatingRequests *mongo.Collection

	colValidatedRequests *mongo.Collection

	colUnknownRequests *mongo.Collection
)

// Initialize the mongodb connection, and store as global variables.
func ConnectToMongoDb(ops *conf.MongoDbOptions) error {
	client, err := mongo.Connect(NewTimoutContext(10), options.Client().ApplyURI(ops.Uri))
	if err != nil {
		return err
	}
	err = client.Ping(NewTimoutContext(2), readpref.Primary())
	if err != nil {
		return err
	}
	app = ops
	mDbClient = client
	colCrawlerRequests = GetColRequests(conf.ColCrawlerRequests)
	colGeneralRequests = GetColRequests(conf.ColGeneralRequests)
	colLandingRequests = GetColRequests(conf.ColLandingRequests)
	colValidatingRequests = GetColRequests(conf.ColValidatingRequests)
	colValidatedRequests = GetColRequests(conf.ColValidatedRequests)
	colUnknownRequests = GetColRequests(conf.ColUnknownRequests)
	return nil
}

func GetColRequests(colName string) *mongo.Collection {
	return mDbClient.Database(app.DbName).Collection(app.GetColName(colName))
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
		raw, _err := colGeneralRequests.FindOne(newCrudContext(), bson.M{}).DecodeBytes()
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

// Calling this func asynchronous is recommended.
func InsertRecordWithErrorProcessed(record *Record) {
	err := InsertRecord(record)
	if err != nil {
		recorder.PrintFailedRecord(record)
	}
}
