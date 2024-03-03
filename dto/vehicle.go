package dto

import (
	"github.com/google/uuid"
)

type VehicleCreationInput struct {
	OwnerID   *uuid.UUID `json:"ownerId" form:"ownerId" binding:""`
	CompanyID *uuid.UUID `json:"companyId" form:"companyId" binding:""`
	//
	ModelID *uuid.UUID `json:"modelId" form:"modelId" binding:""`
	//
	VinNumber *string `json:"vin,omitempty"`
	MakeName  *string `json:"makeName,omitempty"`
	ModelName *string `json:"modelName,omitempty"`
	//
	Color              *string `json:"color" form:"color" binding:"required"`
	Published          bool    `json:"published" form:"published" binding:""`
	Description        string  `json:"description" form:"description" binding:""`
	LicensePlateNumber *string `json:"licensePlateNumber" form:"licensePlateNumber" binding:"required"`
}

type VehicleUpdateInput struct {
	ModelID            *uuid.UUID `json:"modelId" form:"modelId" binding:""`
	OwnerID            *uuid.UUID `json:"ownerId" form:"ownerId" binding:""`
	Published          *bool      `json:"published" form:"published" binding:""`
	Description        *string    `json:"description" form:"description" binding:""`
	LicensePlateNumber *string    `json:"licensePlateNumber" form:"licensePlateNumber" binding:""`

	VinNumber *string `json:"vin,omitempty"`
	Color     *string `json:"color" form:"color" binding:""`
	MakeName  *string `json:"makeName,omitempty"`
	ModelName *string `json:"modelName,omitempty"`
}
