package db

import (
	"github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/recorder"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	app *conf.MongoDbOptions

	mDbClient *mongo.Client

	col *mongo.Collection
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
	col = GetColRequests()
	return nil
}

func GetColRequests() *mongo.Collection {
	return mDbClient.Database(app.DbName).Collection(app.GetColName(conf.ColRequests))
}

func InsertRecord(record *Record) error {
	_, err := col.InsertOne(newCrudContext(), record)
	return err
}

// Calling this func asynchronous is recommended.
func InsertRecordWithErrorProcessed(record *Record) {
	err := InsertRecord(record)
	if err != nil {
		recorder.PrintFailedRecord(record)
	}
}
