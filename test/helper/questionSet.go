package helper

import (
	"context"
	"sync"
	"time"

	faker "github.com/brianvoe/gofakeit/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionSet struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Questions []Question         `bson:"questions"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type Question struct {
	Question string   `bson:"question"`
	Answers  []string `bson:"answers"`
}

func SeedQuestionSetInput(isNew bool) interface{} {
	if isNew {
		return map[string]string{
			"name": faker.Sentence(2),
		}
	} else {
		return QuestionSet{
			Name:      faker.Sentence(2),
			CreatedAt: faker.Date(),
			UpdatedAt: faker.Date(),
			Questions: []Question{
				Question{
					Question: faker.Sentence(2),
					Answers:  []string{faker.Sentence(3), faker.Sentence(3)},
				},
			},
		}
	}
}

func SeedQuestionSet() QuestionSet {
	collection := DB.Collection("questionsets")
	input := SeedQuestionSetInput(false).(QuestionSet)
	insertResult, _ := collection.InsertOne(context.Background(), input)
	input.ID = insertResult.InsertedID.(primitive.ObjectID)

	return input
}

func SeedMultipleQuestionSet(num int) []QuestionSet {
	var listQuestionSet []QuestionSet
	var wg sync.WaitGroup
	c := make(chan QuestionSet, num)

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			c <- SeedQuestionSet()
			wg.Done()
		}()
	}
	wg.Wait()
	close(c)

	for q := range c {
		listQuestionSet = append(listQuestionSet, q)
	}

	return listQuestionSet
}
