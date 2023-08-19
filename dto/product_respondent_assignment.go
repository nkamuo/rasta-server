package dto

import "github.com/google/uuid"

type ProductRespondentAssignmentCreationInput struct {
	ProductID               uuid.UUID `json:"productId" binding:"required"`
	RespondentID            uuid.UUID `json:"respondentId" binding:"required"`
	Note                    string    `json:"note" binding:""`
	Description             string    `json:"description" binding:""`
	Active                  bool      `json:"active" binding:""`
	AllowRespondentActivate bool      `json:"allowRespondentActivate" binding:""`
}

type ProductRespondentAssignmentUpdateInput struct {
	Note                    *string `json:"note" binding:""`
	Description             *string `json:"description" binding:""`
	Active                  *bool   `json:"active" binding:""`
	AllowRespondentActivate *bool   `json:"allowRespondentActivate" binding:""`
}
