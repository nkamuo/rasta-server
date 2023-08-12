package dto

import "github.com/google/uuid"

type RequestInput struct {
	ProductID   *uuid.UUID                      `json:"productId,omitempty"`
	VehicleInfo *RequestVehicleInformationInput `json:"vehicleInfo"`
	FuelInfo    *RequestFuelInformationInput    `json:"fuelInfo"`
	Rate        *uint64                         `json:"rate"`
	Quantity    *uint64                         `json:"quantity,omitempty"`
	Origin      *string                         `json:"origin"`
	Destination *string                         `json:"destination"`
	Description *string                         `json:"description"`
	Note        *string                         `json:"note"`
	// VARIABLE FIELDS - FIELDS THAT MAY BE REQUIRED DEPENDING ON THE TYPE OF REQUEST
	// VehicleModelID     *uuid.UUID `json:"vehicleModelId,omitempty"`
	// FUEL DELIVERY REQUEST
	// FuelTypeCode *string `json:"fuelTypeCode,omitempty"`
}

type RequestFuelInformationInput struct {
	FuelTypeCode        string  `json:"fuelTypeCode,omitempty"`
	FuelTypeDescription *string `json:"fuelTypeDescription"`
}

type RequestVehicleInformationInput struct {
	VehicleModelID     *uuid.UUID `json:"vehicleModelId,omitempty"`
	VehicleDescription *string    `json:"vehicleDescription"`
}
