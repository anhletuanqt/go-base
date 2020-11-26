package questionset

import (
	"base/app/model"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetById(DB *mongo.Database) fiber.Handler {
	collection := DB.Collection("questionsets")
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		var doc model.QuestionSet
		id := c.Params("id")

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.Wrap(err, "")
		}

		filter := bson.D{{"_id", objID}}
		if err := collection.FindOne(ctx, filter).Decode(&doc); err != nil {
			return errors.Wrap(err, "")
		}
		fmt.Println("data", doc)
		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      doc,
		})
	}
}

func GetAll(DB *mongo.Database) fiber.Handler {
	collection := DB.Collection("questionsets")
	ctx := context.Background()
	query := struct {
		Name string `query:"name"`
	}{}

	return func(c *fiber.Ctx) error {
		var docs []model.QuestionSet

		if err := c.QueryParser(&query); err != nil {
			return errors.Wrap(err, "")
		}

		filter := make(map[string]interface{})
		if query.Name != "" {
			// filter["name"] = primitive.Regex{Pattern: query.Name, Options: ""}
			// filter["name"] = bson.D{{"$regex", query.Name}, {"$options", "i"}}
			filter["name"] = map[string]string{
				"$regex":   query.Name,
				"$options": "i",
			}
		}

		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return errors.Wrap(err, "")
		}

		if err := cursor.All(ctx, &docs); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"items":     docs,
		})
	}
}
