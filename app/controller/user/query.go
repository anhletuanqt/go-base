package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type getAllParams struct {
	FullName string `query:"fullName"`
	Email    string `query:"email"`
	Gender   string `query:"gender"`
	Cursor   string `query:"cursor"`
	Limit    int    `query:"limit"`
	Order    string `query:"order"`
}

func GetById(db *mongo.Database) fiber.Handler {
	collection := db.Collection("users")
	FacColl := db.Collection("facilities")
	ctx := context.Background()
	return func(c *fiber.Ctx) error {
		user := make(map[string]interface{})
		id := c.Params("id")

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.Wrap(err, "")
		}

		if err := collection.FindOne(ctx, bson.D{{"_id", objID}}).Decode(&user); err != nil {
			return errors.Wrap(err, "")
		}

		var facility map[string]interface{}
		if fID, ok := user["facility"]; ok {
			if err != nil {
				return errors.Wrap(err, "")
			}
			FacColl.FindOne(ctx, bson.D{{"_id", fID}}).Decode(&facility)
		}
		user["facility"] = facility

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      user,
		})
	}
}

func GetAll(db *mongo.Database) fiber.Handler {
	collection := db.Collection("users")
	ctx := context.Background()
	type objType map[string]interface{}
	return func(c *fiber.Ctx) error {
		users := make([]map[string]interface{}, 0, 10)
		anchor := time.Now()
		queries := getAllParams{}
		if err := c.QueryParser(&queries); err != nil {
			return errors.Wrap(err, "")
		}

		match := objType{
			"$match": buildGetAllQuery(queries),
		}
		lookup := bson.D{{"$lookup", objType{
			"from": "facilities",
			"as":   "facility",
			"let":  objType{"facility": "$facility"},
			"pipeline": []interface{}{
				objType{"$match": objType{
					"$expr": objType{
						"$eq": []interface{}{"$_id", "$$facility"},
					},
				}},
				objType{"$project": objType{
					"createdAt": 0,
					"updatedAt": 0,
				}},
			},
		}}}
		unwind := objType{
			"$unwind": objType{
				"path":                       "$facility",
				"preserveNullAndEmptyArrays": true,
			},
		}
		project := objType{"$project": objType{
			"createdAt": 0,
			"updatedAt": 0,
		}}
		limit := objType{
			"$limit": 1000,
		}
		sort := objType{
			"$sort": objType{
				"email": 1,
			},
		}
		aggregate := []interface{}{
			match,
			lookup,
			unwind,
			project,
			limit,
			sort,
		}
		cursor, err := collection.Aggregate(ctx, aggregate)

		// filter := buildGetAllQuery(queries)
		// opts := buildGetAllOpts(queries)
		// cursor, err := collection.Find(ctx, filter, opts)
		if err != nil {
			return errors.Wrap(err, "")
		}

		if err := cursor.All(ctx, &users); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"items":     users,
			"time":      time.Since(anchor),
		})
	}
}
