package dto

import (
	"github.com/google/uuid"
)

type VehicleCreationInput struct {
	ModelID            uuid.UUID `json:"modelId" binding:"required"`
	OwnerID            uuid.UUID `json:"ownerId" binding:"required"`
	Description        string    `json:"description" binding:""`
	LicensePlaceNumber *string   `json:"licensePlaceNumber" binding:"required"`
}

type VehicleUpdateInput struct {
	ModelID            *uuid.UUID `json:"modelId" binding:"required"`
	OwnerID            *uuid.UUID `json:"ownerId" binding:"required"`
	Description        *string    `json:"description" binding:""`
	LicensePlaceNumber *string    `json:"licensePlaceNumber" binding:"required"`
}
