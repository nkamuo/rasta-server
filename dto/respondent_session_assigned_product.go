package dto

import "github.com/google/uuid"

type RespondentSessionAssignedProductInput struct {
	AssignmentID uuid.UUID `json:"assignmentId" binding:"required"`
	Note         string    `json:"note" binding:""`
	Description  string    `json:"description" binding:""`
	Active       bool      `json:"active" binding:""`
}

type RespondentSessionAssignedProductCreationInput struct {
	RespondentSessionAssignedProductInput
}

type RespondentSessionAssignedProductUpdateInput struct {
	Note        *string `json:"note" binding:""`
	Description *string `json:"description" binding:""`
	Active      *bool   `json:"active" binding:""`
}
