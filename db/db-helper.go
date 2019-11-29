package db

import (
	"github.com/zhanbei/static-server/conf"
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
