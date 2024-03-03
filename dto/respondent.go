package dto

import (
	"github.com/google/uuid"
)

type RespondentCreationInput struct {
	UserID    uuid.UUID  `json:"userId" binding:"required"`
	CompanyID *uuid.UUID `json:"companyId" binding:""`
	PlaceID   *uuid.UUID `json:"placeId" binding:""`
	VehicleID *uuid.UUID `json:"vehicleId" binding:""`
	Active    *bool      `json:"active" binding:""`
	// Title       string `json:"title" binding:"required"`
	// Description string `json:"description" binding:"required"`
}

type RespondentUpdateInput struct {
	CompanyID *uuid.UUID `json:"companyId" binding:""`
	PlaceID   *uuid.UUID `json:"placeId" binding:""`
	VehicleID *uuid.UUID `json:"vehicleId" binding:""`
	Active    *bool      `json:"active" binding:""`
}

type RespondentCompanyAssignmentInput struct {
	RespondentID uuid.UUID `json:"respondentId" binding:"required"`
}

type RespondentDocumentVerificationInput struct {
	Ssn *string `form:"ssn" json:"ssn" binding:""`
	// Documents *[]*ImageDocument `json:"documents,omitempty"`
	// FileData *[]*multipart.FileHeader `form:"documents" binding:""`
}
