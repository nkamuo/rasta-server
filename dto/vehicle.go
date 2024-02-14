package dto

import (
	"github.com/google/uuid"
)

type VehicleCreationInput struct {
	OwnerID   *uuid.UUID `json:"ownerId" binding:""`
	CompanyID *uuid.UUID `json:"companyId" binding:""`
	//
	ModelID *uuid.UUID `json:"modelId" binding:""`
	//
	VinNumber *string `json:"vin,omitempty"`
	MakeName  *string `json:"makeName,omitempty"`
	ModelName *string `json:"modelName,omitempty"`
	//
	Color              *string `json:"color" binding:"required"`
	Published          bool    `json:"published" binding:""`
	Description        string  `json:"description" binding:""`
	LicensePlateNumber *string `json:"licensePlateNumber" binding:"required"`
}

type VehicleUpdateInput struct {
	ModelID            *uuid.UUID `json:"modelId" binding:""`
	OwnerID            *uuid.UUID `json:"ownerId" binding:""`
	Published          *bool      `json:"published" binding:""`
	Description        *string    `json:"description" binding:""`
	LicensePlateNumber *string    `json:"licensePlateNumber" binding:""`

	VinNumber *string `json:"vin,omitempty"`
	Color     *string `json:"color" binding:""`
	MakeName  *string `json:"makeName,omitempty"`
	ModelName *string `json:"modelName,omitempty"`
}
