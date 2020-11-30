package facilities

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetById(db *mongo.Database) fiber.Handler {
	collection := db.Collection("facilities")
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		facility := map[string]interface{}{}

		objID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return errors.Wrap(err, "")
		}

		query := bson.D{{"_id", objID}}
		opts := options.FindOne()
		opts.SetProjection(map[string]interface{}{
			"createdAt": 0,
			"updatedAt": 0,
		})
		if err := collection.FindOne(ctx, query, opts).Decode(&facility); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      facility,
		})
	}
}

func GetAll(db *mongo.Database) fiber.Handler {
	collection := db.Collection("facilities")
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		facilities := make([]map[string]interface{}, 0, 10)

		cursor, err := collection.Find(ctx, bson.D{{}})
		if err != nil {
			return errors.Wrap(err, "")
		}

		if err := cursor.All(ctx, &facilities); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"items":     facilities,
		})
	}
}
