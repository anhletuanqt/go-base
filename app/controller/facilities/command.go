package facilities

import (
	"base/app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func Create(db *mongo.Database) fiber.Handler {
	collection := db.Collection("facilities")
	validate := validator.New()
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		facility := model.Facility{}

		if err := c.BodyParser(&facility); err != nil {
			return errors.Wrap(err, "")
		}

		if err := validate.Struct(facility); err != nil {
			return errors.Wrap(err, "")
		}

		facility.BeforeUpdate(true)

		insertedDoc, err := collection.InsertOne(ctx, facility)
		if err != nil {
			return errors.Wrap(err, "")
		}

		facility.ID = insertedDoc.InsertedID.(primitive.ObjectID)

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      facility,
		})
	}
}

func Delete(db *mongo.Database) fiber.Handler {
	collection := db.Collection("facilities")
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		response := make(map[string]interface{})

		objID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return errors.Wrap(err, "")
		}

		if err := collection.FindOneAndDelete(ctx, bson.D{{"_id", objID}}).Decode(&response); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      response,
		})
	}
}

func Update(db *mongo.Database) fiber.Handler {
	collection := db.Collection("facilities")
	ctx := context.Background()
	validate := validator.New()

	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		response := make(map[string]interface{})

		objID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return errors.Wrap(err, "")
		}

		input := model.UpdateFacility{}
		if err := c.BodyParser(&input); err != nil {
			return errors.Wrap(err, "")
		}

		if err := validate.Struct(input); err != nil {
			return errors.Wrap(err, "")
		}

		input.UpdatedAt = time.Now()
		updateData := map[string]interface{}{
			"$set": input,
		}

		opts := options.FindOneAndUpdate()
		opts.SetReturnDocument(1)

		if err := collection.FindOneAndUpdate(ctx, bson.D{{"_id", objID}}, updateData, opts).Decode(&response); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      response,
		})
	}
}
