package helper

import (
	"context"
	"sync"
	"time"

	faker "github.com/brianvoe/gofakeit/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Facility struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Admin     *primitive.ObjectID `bson:"admin,omitempty"`
	Name      string              `bson:"name"`
	Phone     string              `bson:"phone"`
	Status    string              `bson:"status"`
	Locations []string            `bson:"locations,omitempty"`
	CreatedAt time.Time           `bson:"createdAt,omitempty"`
	UpdatedAt time.Time           `bson:"updatedAt,omitempty"`
}

func seedStatus(status ...string) string {
	if len(status) != 0 {
		return status[0]
	}

	i := faker.RandomInt([]int{0, 1})
	if i == 0 {
		return "Close"
	}

	return "Active"
}

func SeedFacilityInput() map[string]interface{} {
	input := make(map[string]interface{})

	input["name"] = faker.Name()
	input["phone"] = "+" + faker.Phone()
	input["status"] = seedStatus()
	input["locations"] = []string{
		faker.City(),
		faker.City(),
	}

	return input
}

func SeedFacility() Facility {
	input := SeedFacilityInput()
	facility := Facility{
		CreatedAt: faker.Date(),
		UpdatedAt: faker.Date(),
		Name:      input["name"].(string),
		Phone:     input["phone"].(string),
		Status:    input["status"].(string),
		Locations: input["locations"].([]string),
	}
	collection := DB.Collection("facilities")
	result, _ := collection.InsertOne(context.Background(), facility)
	facility.ID = result.InsertedID.(primitive.ObjectID)

	return facility
}

func SeedMultipleFacility(num int) []Facility {
	c := make(chan Facility, num)
	facilities := make([]Facility, 0, num)

	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			c <- SeedFacility()
			wg.Done()
		}()
	}
	wg.Wait()
	close(c)

	for v := range c {
		facilities = append(facilities, v)
	}

	return facilities
}
