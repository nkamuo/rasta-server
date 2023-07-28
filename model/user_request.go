package model

import (
	"github.com/google/uuid"
)

type UserRequest struct {
	// gorm.Model;
	ID   uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	user User      `` // THE User MAKING THE REQUEST
}

type VehicleTowUserRequest struct {
	UserRequest
	origin      UserRequestPosition
	destination UserRequestPosition
}

type UserRequestPosition struct {
	coordinate UserRequestPositionCoordinate
}

type UserRequestPositionCoordinate struct {
	latitude  float32
	longitude float32
	altitude  float32
}

func (UserRequest) TableName() string {
	return "user_requests"
}

func (VehicleTowUserRequest) TableName() string {
	return "user_requests"
}
