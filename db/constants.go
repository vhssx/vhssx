package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type ObjectId = primitive.ObjectID

var NewObjectId = primitive.NewObjectID
