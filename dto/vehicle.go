package dto

import (
	"github.com/google/uuid"
)

type VehicleCreationInput struct {
	ModelID            uuid.UUID `json:"modelId" binding:"required"`
	OwnerID            uuid.UUID `json:"ownerId" binding:"required"`
	Color              *string   `json:"color" binding:"required"`
	Published          bool      `json:"published" binding:""`
	Description        string    `json:"description" binding:""`
	LicensePlateNumber *string   `json:"licensePlateNumber" binding:"required"`
}

type VehicleUpdateInput struct {
	ModelID            *uuid.UUID `json:"modelId" binding:""`
	OwnerID            *uuid.UUID `json:"ownerId" binding:""`
	Color              *string    `json:"color" binding:""`
	Published          *bool      `json:"published" binding:""`
	Description        *string    `json:"description" binding:""`
	LicensePlateNumber *string    `json:"licensePlateNumber" binding:""`
}
