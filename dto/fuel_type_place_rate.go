package dto

import "github.com/google/uuid"

type FuelTypePlaceRateCreationInput struct {
	Rate uint64 `json:"rate" binding:"required"`
	//
	FuelTypeID uuid.UUID `json:"fuelTypeId,omitempty" binding:"required"`
	PlaceID    uuid.UUID `gorm:"not null" json:"placeId,omitempty" binding:"required"`
	//
	Description *string `json:"description,omitempty" binding:""`
	Active      bool    `json:"active" binding:""`
}

type FuelTypePlaceRateUpdateInput struct {
	Rate        *uint64 `json:"rate" binding:""`
	Active      *bool   `json:"active" binding:""`
	Description *string `json:"description,omitempty" binding:""`
}
