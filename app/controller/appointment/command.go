package appointment

import (
	"base/app/model"
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func Create(db *mongo.Database) fiber.Handler {
	collection := db.Collection("appointments")
	ctx := context.Background()
	validate := validator.New()

	return func(c *fiber.Ctx) error {
		appointment := model.Appointment{}

		if err := c.BodyParser(&appointment); err != nil {
			return errors.Wrap(err, "")
		}

		if err := validate.Struct(appointment); err != nil {
			return errors.Wrap(err, "")
		}

		appointment.BeforeSave()

		insertResult, err := collection.InsertOne(ctx, appointment)
		if err != nil {
			return errors.Wrap(err, "")
		}

		InsertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
		if !ok {
			return errors.Wrap(errors.New("InsertedID is not ObjectID"), "")
		}
		appointment.ID = InsertedID

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"message":   "Create an appointment successfully",
			"item":      appointment,
		})
	}
}

func Delete(db *mongo.Database) fiber.Handler {
	collection := db.Collection("appointments")
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		appointment := make(map[string]interface{})

		id := c.Params("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.Wrap(err, "")
		}

		filter := bson.D{{"_id", objID}}
		if err := collection.FindOneAndDelete(ctx, filter).Decode(&appointment); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      appointment,
		})
	}
}

func Update(db *mongo.Database) fiber.Handler {
	collection := db.Collection("appointments")
	ctx := context.Background()
	validate := validator.New()

	return func(c *fiber.Ctx) error {
		appointment := make(map[string]interface{})
		updateInput := model.UpdatedAppointment{}

		id := c.Params("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.Wrap(err, "")
		}

		if err := c.BodyParser(&updateInput); err != nil {
			return errors.Wrap(err, "")
		}

		if err := validate.Struct(updateInput); err != nil {
			return errors.Wrap(err, "")
		}

		filter := bson.D{{"_id", objID}}
		opts := options.FindOneAndUpdate()
		opts.SetReturnDocument(1)
		if err := collection.FindOneAndUpdate(ctx, filter, bson.D{{"$set", updateInput}}, opts).Decode(&appointment); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      appointment,
		})
	}
}
