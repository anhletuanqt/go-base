package user

import (
	"base/app/model"
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Create(db *mongo.Database) fiber.Handler {
	collection := db.Collection("users")
	ctx := context.Background()
	validate := validator.New()

	return func(c *fiber.Ctx) error {
		user := model.User{}

		if err := c.BodyParser(&user); err != nil {
			return errors.Wrap(err, "")
		}

		if err := validate.Struct(user); err != nil {
			return errors.Wrap(err, "")
		}

		user.BeforeSave()

		insertedResult, err := collection.InsertOne(ctx, user)
		if err != nil {
			return errors.Wrap(err, "")
		}
		user.ID = insertedResult.InsertedID.(primitive.ObjectID)

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      user,
		})
	}
}

func Delete(db *mongo.Database) fiber.Handler {
	collection := db.Collection("users")
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		user := make(map[string]interface{})
		id := c.Params("id")

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.Wrap(err, "")
		}

		if err := collection.FindOneAndDelete(ctx, bson.D{{"_id", objID}}).Decode(&user); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      user,
		})
	}
}

func Update(db *mongo.Database) fiber.Handler {
	collection := db.Collection("users")
	ctx := context.Background()
	validate := validator.New()

	return func(c *fiber.Ctx) error {
		user := make(map[string]interface{})
		input := model.UpdateUser{}
		id := c.Params("id")

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.Wrap(err, "")
		}

		if err := c.BodyParser(&input); err != nil {
			return errors.Wrap(err, "")
		}

		if err := validate.Struct(input); err != nil {
			return errors.Wrap(err, "")
		}

		input.BeforeUpdate()
		updateData := map[string]interface{}{
			"$set": input,
		}

		opts := options.FindOneAndUpdate()
		opts.SetReturnDocument(1)
		if err := collection.FindOneAndUpdate(ctx, bson.D{{"_id", objID}}, updateData, opts).Decode(&user); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      user,
		})
	}
}
