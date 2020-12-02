package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	FullName  string              `bson:"fullName" json:"fullName" validate:"required"`
	Gender    string              `bson:"gender" json:"gender" validate:"oneof=Male Female"`
	Dob       string              `bson:"dob" json:"dob" validate:"required"`
	Email     string              `bson:"email" json:"email" validate:"required,email"`
	Types     []string            `bson:"types" json:"types" validate:"dive,oneof=Patient Physician Technician 'Facility Admin'"`
	CreatedAt time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time           `baon:"updatedAt" json:"updatedAt"`
	Facility  *primitive.ObjectID `bson:"facility,omitempty" json:"facility,omitempty"`
}

func (u *User) BeforeSave() {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

type UpdateUser struct {
	FullName  string              `bson:"fullName" json:"fullName,omitempty" validate:"required"`
	Gender    string              `bson:"gender" json:"gender,omitempty" validate:"oneof=Male Female"`
	Dob       string              `bson:"dob" json:"dob,omitempty" validate:"required"`
	Email     string              `bson:"email" json:"email,omitempty" validate:"required,email"`
	Types     []string            `bson:"types,omitempty" json:"types,omitempty" validate:"dive,oneof=Patient Physician Technician 'Facility Admin'"`
	Facility  *primitive.ObjectID `bson:"facility,omitempty" json:"facility,omitempty"`
	UpdatedAt time.Time           `baon:"updatedAt" json:"updatedAt"`
}

func (u *UpdateUser) BeforeUpdate() {
	u.UpdatedAt = time.Now()
}
