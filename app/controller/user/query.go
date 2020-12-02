package user

import (
	"context"

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

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      user,
		})
	}
}

func GetAll(db *mongo.Database) fiber.Handler {
	collection := db.Collection("users")
	ctx := context.Background()
	return func(c *fiber.Ctx) error {
		users := make([]map[string]interface{}, 0, 10)

		queries := getAllParams{}
		if err := c.QueryParser(&queries); err != nil {
			return errors.Wrap(err, "")
		}

		filter := buildGetAllQuery(queries)
		opts := buildGetAllOpts(queries)
		cursor, err := collection.Find(ctx, filter, opts)
		if err != nil {
			return errors.Wrap(err, "")
		}

		if err := cursor.All(ctx, &users); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"items":     users,
		})
	}
}
