package database

import (
	"base/config"
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func Connect(conf *config.Config) {
	clientOptions := options.Client().ApplyURI(conf.Mongo.URL)
	// Connect to MongoDB
	withTimeOut, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(withTimeOut, clientOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db = client.Database(conf.Mongo.DBName)

	// Index for question set
	err = createQuestionSetIndex(db)
	if err != nil {
		fmt.Println("Question set index err:", err)
	}

	// Index for facility
	err = createFacilityIndex(db)
	if err != nil {
		fmt.Println("Facility index err:", err)
	}

	// Index for user
	err = createUserIndex(db)
	if err != nil {
		fmt.Println("User index err:", err)
	}
}

func GetDB() *mongo.Database {
	return db
}

func createQuestionSetIndex(db *mongo.Database) error {
	coll := db.Collection("questionsets")

	// 1. Lets define the keys for the index we want to create
	mod := []mongo.IndexModel{
		mongo.IndexModel{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := coll.Indexes().CreateMany(context.Background(), mod)

	return err
}

func createFacilityIndex(db *mongo.Database) error {
	coll := db.Collection("facilities")

	// 1. Lets define the keys for the index we want to create
	mod := []mongo.IndexModel{
		mongo.IndexModel{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := coll.Indexes().CreateMany(context.Background(), mod)

	return err
}

func createUserIndex(db *mongo.Database) error {
	coll := db.Collection("users")

	// 1. Lets define the keys for the index we want to create
	mod := []mongo.IndexModel{
		mongo.IndexModel{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := coll.Indexes().CreateMany(context.Background(), mod)

	return err
}
