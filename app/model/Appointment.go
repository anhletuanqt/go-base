package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Patient         *primitive.ObjectID  `bson:"patient" json:"patient"`
	Technician      *primitive.ObjectID  `bson:"technician" json:"technician"`
	Physician       *primitive.ObjectID  `bson:"physician" json:"physician"`
	Facility        *primitive.ObjectID  `bson:"facility" json:"facility" validate:"required"`
	CanFinishedUser []primitive.ObjectID `bson:"canFinishedUser" json:"canFinishedUser"`
	Status          string               `bson:"status" json:"status"`
	MeetingAt       time.Time            `bson:"meetingAt" json:"meetingAt" validate:"required"`
	CreatedAt       time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time            `bson:"updatedAt" json:"updatedAt"`
}

func (a *Appointment) BeforeSave() {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	if a.Technician != nil {
		a.CanFinishedUser = append(a.CanFinishedUser, *a.Technician)
	}

	if a.Physician != nil {
		a.CanFinishedUser = append(a.CanFinishedUser, *a.Physician)
	}

	if a.Status == "" {
		a.Status = "Pending"
	}
}

type UpdatedAppointment struct {
	Technician *primitive.ObjectID `bson:"technician,omitempty" json:"technician"`
	Physician  *primitive.ObjectID `bson:"physician,omitempty" json:"physician"`
	Status     string              `bson:"status,omitempty" json:"status" validate:"omitempty,oneof=Approved Canceled Overdue Finish"`
	MeetingAt  time.Time           `bson:"meetingAt,omitempty" json:"meetingAt"`
	UpdatedAt  time.Time           `bson:"updatedAt,omitempty" json:"updatedAt"`
}

func (u *UpdatedAppointment) BeforeUpdate() {
	u.UpdatedAt = time.Now()
}
