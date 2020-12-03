package helper

import (
	"context"
	"sync"
	"time"

	faker "github.com/brianvoe/gofakeit/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	FullName  string              `bson:"fullName"`
	Gender    string              `bson:"gender"`
	Dob       string              `bson:"dob"`
	Email     string              `bson:"email"`
	Types     []string            `bson:"types"`
	CreatedAt time.Time           `bson:"createdAt"`
	UpdatedAt time.Time           `baon:"updatedAt"`
	Facility  *primitive.ObjectID `bson:"facility,omitempty"`
}

func seedGender() string {
	i := faker.RandomInt([]int{0, 1})

	if i == 1 {
		return "Male"
	}

	return "Female"
}

func seedTypes() []string {
	i := faker.RandomInt([]int{0, 1, 2, 3})
	types := make([]string, 1)

	switch i {
	case 0:
		types[0] = "Patient"
	case 1:
		types[0] = "Technician"
	case 2:
		types[0] = "Physician"
	case 3:
		types[0] = "Facility Admin"
	}

	return types
}

func SeedUserInput(facility ...primitive.ObjectID) map[string]interface{} {
	input := make(map[string]interface{})
	input["fullName"] = faker.Name()
	input["gender"] = seedGender()
	input["dob"] = faker.Date().String()
	input["email"] = faker.Email()
	input["types"] = seedTypes()
	if len(facility) != 0 {
		input["facility"] = &facility[0]
	}

	return input
}

func SeedUser() User {
	collection := DB.Collection("users")
	input := SeedUserInput()
	user := User{
		FullName:  input["fullName"].(string),
		Gender:    input["gender"].(string),
		Dob:       input["dob"].(string),
		Email:     input["email"].(string),
		Types:     input["types"].([]string),
		CreatedAt: faker.Date(),
		UpdatedAt: faker.Date(),
	}

	facility, ok := input["facility"]
	if ok {
		user.Facility = facility.(*primitive.ObjectID)
	}

	inserted, _ := collection.InsertOne(context.Background(), user)
	user.ID = inserted.InsertedID.(primitive.ObjectID)

	return user
}

func SeedMultipleUser(num int) []User {
	users := make([]User, 0, num)
	c := make(chan User, num)
	var wg sync.WaitGroup

	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			c <- SeedUser()
			wg.Done()
		}()
	}
	wg.Wait()
	close(c)

	for u := range c {
		users = append(users, u)
	}

	return users
}
