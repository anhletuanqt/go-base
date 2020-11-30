package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Facility struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"_id"`
	Admin     *primitive.ObjectID `bson:"admin,omitempty" json:"admin,omitempty"`
	Name      string              `bson:"name" json:"name" validate:"required"`
	Phone     string              `bson:"phone" json:"phone" validate:"required,e164"`
	Status    string              `bson:"status" json:"status" validate:"required,oneof=Active Close"`
	Locations []string            `bson:"locations,omitempty" json:"locations,omitempty" validate:"required"`
	CreatedAt time.Time           `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time           `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

func (f *Facility) BeforeUpdate(isNew bool) {
	if isNew {
		f.CreatedAt = time.Now()
		f.UpdatedAt = time.Now()
	} else {
		f.UpdatedAt = time.Now()
	}
}

type UpdateFacility struct {
	Admin     *primitive.ObjectID `bson:"admin,omitempty" json:"admin,omitempty"`
	Name      string              `bson:"name,omitempty" json:"name,omitempty"`
	Phone     string              `bson:"phone" json:"phone" validate:"e164"`
	Status    string              `bson:"status,omitempty" json:"status,omitempty" validate:"omitempty,oneof=Active Close"`
	Locations []string            `bson:"locations,omitempty" json:"locations,omitempty"`
	UpdatedAt time.Time           `bson:"updatedAt" json:"updatedAt"`
}
