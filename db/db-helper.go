package db

import (
	"github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/recorder"
	"github.com/zhanbei/static-server/secoo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	app *conf.MongoDbOptions

	mDbClient *mongo.Client

	colCrawlerRequests *mongo.Collection

	colGeneralRequests *mongo.Collection

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
	case secoo.LevelFollowingTimeRequest:
		// ID = NewID(); SID = NextSessionID;
		record.SessionId = &record.Id
		record.Id = NewObjectId()
		record.Session = nil
		_, err = colValidatedRequests.InsertOne(newCrudContext(), record)
	default:
		// ID = NewID(); SID = nil;
		record.SessionId = &record.Id
		// Create a new ID for whatever safety reason.
		record.Id = NewObjectId()
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
