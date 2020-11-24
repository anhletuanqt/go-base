package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionSet struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name      string             `json:"name" validate:"required"`
	Questions []Question         `json:"questions"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Question struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Question  string             `json:"question"`
	Answers   []string           `json:"answers"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type UpdateQuestionSet struct {
	Name      string     `json:"name"`
	Questions []Question `json:"questions"`
	UpdatedAt time.Time  `json:"updatedAt" bson:"updatedAt"`
}

func NewQuestionSet() *QuestionSet {
	return &QuestionSet{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Questions: []Question{},
	}
}
