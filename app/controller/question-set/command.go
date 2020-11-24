package questionset

import (
	"base/app/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(DB *mongo.Database) fiber.Handler {
	collection := DB.Collection("questionsets")
	ctx := context.Background()
	validate := validator.New()

	return func(c *fiber.Ctx) error {
		newDoc := model.NewQuestionSet()

		if err := c.BodyParser(newDoc); err != nil {
			return errors.Wrap(err, "")
		}

		if err := validate.Struct(newDoc); err != nil {
			return errors.Wrap(err, "")
		}

		inserted, err := collection.InsertOne(ctx, newDoc)
		if err != nil {
			return errors.Wrap(err, "")
		}

		objID, ok := inserted.InsertedID.(primitive.ObjectID)
		if !ok {
			return errors.Wrap(errors.New("This isn't ObjectID"), "")
		}
		newDoc.ID = objID

		return c.JSON(fiber.Map{
			"isSuccess": "ok",
			"item":      newDoc,
		})
	}
}

func Delete(DB *mongo.Database) fiber.Handler {
	collection := DB.Collection("questionsets")
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		var deletedDoc model.QuestionSet
		id := c.Params("id")

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.Wrap(err, "")
		}

		filter := bson.D{{"_id", objID}}
		if err := collection.FindOneAndDelete(ctx, filter).Decode(&deletedDoc); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": "ok",
			"item":      deletedDoc,
		})
	}
}

func UpdateById(DB *mongo.Database) fiber.Handler {
	collection := DB.Collection("questionsets")
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		var updatedDoc model.QuestionSet
		updatedInput := model.UpdateQuestionSet{
			UpdatedAt: time.Now(),
			Questions: []model.Question{},
		}

		if err := c.BodyParser(&updatedInput); err != nil {
			return errors.Wrap(err, "")
		}

		id := c.Params("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}

		updatedDate := map[string]interface{}{
			"$set": updatedInput,
		}
		opts := options.FindOneAndUpdate()
		opts.SetReturnDocument(1)
		if err := collection.FindOneAndUpdate(ctx, bson.D{{"_id", objID}}, updatedDate, opts).Decode(&updatedDoc); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      updatedDoc,
		})
	}
}
