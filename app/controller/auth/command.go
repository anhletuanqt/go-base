package auth

import (
	"base/utils"
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(db *mongo.Database) fiber.Handler {
	collection := db.Collection("users")
	ctx := context.Background()
	type LoginBody struct {
		ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Email    string             `json:"email" bson:"email" validate:"required,email"`
		Password string             `json:"password" bson:"password" validate:"required"`
	}
	validate := validator.New()

	return func(c *fiber.Ctx) error {
		var user LoginBody
		if err := c.BodyParser(&user); err != nil {
			return errors.Wrap(err, "")
		}

		if err := validate.Struct(user); err != nil {
			return errors.Wrap(err, "")
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.Wrap(err, "")
		}
		user.Password = string(hashPassword)

		insertResult, err := collection.InsertOne(ctx, user)
		if err != nil {
			return errors.Wrap(err, "")
		}

		insertObjID, ok := insertResult.InsertedID.(primitive.ObjectID)
		if !ok {
			return errors.Wrap(errors.New("Asserting error"), "")
		}
		user.ID = insertObjID

		signData := map[string]interface{}{
			"_id":   user.ID,
			"email": user.Email,
		}
		token, err := utils.SignToken(signData)
		if err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"token":     token,
		})
	}
}

func SignIn(db *mongo.Database) fiber.Handler {
	type Login struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	validate := validator.New()
	collection := db.Collection("users")
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		loginData := Login{}
		user := make(map[string]interface{})

		if err := c.BodyParser(&loginData); err != nil {
			return errors.Wrap(err, "")
		}

		if err := validate.Struct(loginData); err != nil {
			return errors.Wrap(err, "")
		}

		filter := bson.D{{"email", loginData.Email}}
		opts := options.FindOne()
		opts.SetProjection(map[string]interface{}{
			"email":    1,
			"password": 1,
		})
		if err := collection.FindOne(ctx, filter, opts).Decode(&user); err != nil {
			return errors.Wrap(err, "")
		}

		// Check password
		if err := bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(loginData.Password)); err != nil {
			return errors.Wrap(err, "")
		}

		delete(user, "password")
		token, err := utils.SignToken(user)
		if err != nil {
			return errors.Wrap(err, "")
		}

		return c.JSON(fiber.Map{
			"isSuccess": true,
			"token":     token,
		})
	}
}
