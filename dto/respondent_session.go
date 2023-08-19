package dto

import "github.com/google/uuid"

type RespondentSessionCreationInput struct {
	// ProductID           uuid.UUID                `json:"productId" binding:"required"`
	RespondentID        uuid.UUID                               `json:"respondentId" binding:"required"`
	StartingCoordinates LocationCoordinatesInput                `json:"startingCoords" binding:"required"`
	Assignments         []RespondentSessionAssignedProductInput `json:"products" binding:"required"`
	Note                string                                  `json:"note" binding:""`
	Description         string                                  `json:"description" binding:""`
	Active              bool                                    `json:"active" binding:""`
}

type RespondentSessionUpdateInput struct {
	Note        *string `json:"note" binding:""`
	Description *string `json:"description" binding:""`
	Active      *bool   `json:"active" binding:""`
}
