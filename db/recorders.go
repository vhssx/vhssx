package db

import (
	"github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/recorder"
)

// Get a default recorder if there are no records at all.
func GetMongoRecorder(ops *conf.MongoDbOptions) recorder.IRecorder {
	if ops == nil {
		return nil
	}
	return NewRecorder(ops)
}
