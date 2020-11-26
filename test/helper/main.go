package helper

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Database

func Init(mDB *mongo.Database) {
	DB = mDB
}
