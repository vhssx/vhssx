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

	colGeneralRequests *mongo.Collection

	colValidatingRequests *mongo.Collection

	colValidatedRequests *mongo.Collection
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
	colGeneralRequests = GetColRequests(conf.ColRequests)
	colValidatingRequests = GetColRequests(conf.ColValidatingRequests)
	colValidatedRequests = GetColRequests(conf.ColValidatedRequests)
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
	switch record.Session.Level {
	case secoo.LevelSecondTimeRequest:
		record.Session = nil
		_, err = colValidatingRequests.InsertOne(newCrudContext(), record)
	case secoo.LevelFollowingTimeRequest:
		record.Session = nil
		_, err = colValidatedRequests.InsertOne(newCrudContext(), record)
	case secoo.LevelFirstTimeRequest:
		fallthrough
	default:
		_, err = colGeneralRequests.InsertOne(newCrudContext(), record)
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
