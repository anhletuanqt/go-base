package appointment

import (
	"base/app/model"
	"base/database"
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
		authUser := c.Locals("user").(map[string]interface{})
		appointment := model.Appointment{}

		if err := c.BodyParser(&appointment); err != nil {
			return errors.Wrap(err, "")
		}

		if err := validate.Struct(appointment); err != nil {
			return errors.Wrap(err, "")
		}

		patientID := authUser["_id"].(string)
		patientObjID, err := primitive.ObjectIDFromHex(patientID)
		if err != nil {
			return errors.Wrap(err, "")
		}
		appointment.Patient = &patientObjID
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

func Approve(db *mongo.Database) fiber.Handler {
	collection := db.Collection("appointments")
	userColl := db.Collection("users")
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		var err error
		var session mongo.Session

		authUser := c.Locals("user").(map[string]interface{})
		client := database.GetClient()
		appointment := make(map[string]interface{})

		id := c.Params("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.Wrap(err, "")
		}

		// Start session
		if session, err = client.StartSession(); err != nil {
			return errors.Wrap(err, "")
		}
		defer session.EndSession(ctx)

		// Start transaction
		if err = session.StartTransaction(); err != nil {
			return errors.Wrap(err, "")
		}

		callback := func(sc mongo.SessionContext) error {
			user := make(map[string]interface{})

			// Find appointment
			filter := bson.D{{"_id", objID}}
			update := bson.D{{"status", "Approved"}}
			opts := options.FindOneAndUpdate()
			opts.SetReturnDocument(1)
			if err = collection.FindOneAndUpdate(sc, filter, update, opts).Decode(&appointment); err != nil {
				return errors.Wrap(err, "")
			}

			// Find admin user
			userId := authUser["_id"].(string)
			userObjID, err := primitive.ObjectIDFromHex(userId)
			userOpts := options.FindOne()
			userOpts.SetProjection(map[string]interface{}{
				"facility": 1,
				"types":    1,
			})
			if err != nil {
				return errors.Wrap(err, "")
			}

			if err := userColl.FindOne(sc, bson.D{{"_id", userObjID}}, userOpts).Decode(&user); err != nil {
				return errors.Wrap(err, "")
			}

			facilityID := appointment["facility"].(string)
			isAdmin := checkFacilityAdmin(map[string]interface{}{
				"facility": user["facility"],
				"types":    user["types"],
			}, facilityID)

			if !isAdmin {
				return errors.Wrap(errors.New("Auth user is not Facility Admin"), "")
			}

			if err := session.CommitTransaction(sc); err != nil {
				return errors.Wrap(err, "")
			}

			return nil
		}

		if err = mongo.WithSession(ctx, session, callback); err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"item":      "item",
		})
	}
}
