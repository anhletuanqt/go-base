package test

import (
	"base/app"
	"base/config"
	"base/database"
	"base/server"
	"base/test/helper"
	"context"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var DB *mongo.Database
var App *fiber.App

const BASE_URL string = "http://localhost:3000/api/v1"

var _ = BeforeSuite(func() {
	ctx := context.Background()
	conf := config.New()
	database.Connect(conf)
	DB = database.GetDB()

	listCollections, _ := DB.ListCollectionNames(ctx, bson.D{{}})
	for _, v := range listCollections {
		DB.Collection(v).Drop(ctx)
	}
	helper.Init(DB)
	App = server.Setup()
	app.InitRoute(App)
})
