package dto

import "github.com/google/uuid"

type RespondentCreationInput struct {
	UserId uuid.UUID `json:"userId" binding:"required"`
	// Title       string `json:"title" binding:"required"`
	// Description string `json:"description" binding:"required"`
}

type RespondentUpdateInput struct {
}

type RespondentCompanyAssignmentInput struct {
	RespondentID uuid.UUID `json:"respondentId" binding:"required"`
}
