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

	colShortenerRedirections *mongo.Collection
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
	colCrawlerRequests = getCol(conf.ColCrawlerRequests)
	colGeneralRequests = getCol(conf.ColGeneralRequests)
	colLandingRequests = getCol(conf.ColLandingRequests)
	colValidatingRequests = getCol(conf.ColValidatingRequests)
	colValidatedRequests = getCol(conf.ColValidatedRequests)
	colUnknownRequests = getCol(conf.ColUnknownRequests)
	colShortenerRedirections = getCol(conf.ColShortenerRedirections)
	return nil
}

func GetCol(colName string) *mongo.Collection {
	return mDbClient.Database(app.DbName).Collection(colName)
}

func getCol(colName string) *mongo.Collection {
	return GetCol(app.GetColName(colName))
}
