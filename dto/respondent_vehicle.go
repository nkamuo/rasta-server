package dto

import "github.com/google/uuid"

type RespondentVehicleSelectionInput struct {
	// ProductID           uuid.UUID                `json:"productId" binding:"required"`
	VehicleID   uuid.UUID `json:"vehicleId" binding:"required"`
	Note        string    `json:"note" binding:""`
	Description string    `json:"description" binding:""`
	Active      *bool     `json:"active" binding:""`
}
