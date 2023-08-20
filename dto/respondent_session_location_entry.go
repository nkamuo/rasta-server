package dto

// import "github.com/google/uuid"

type RespondentSessionLocationEntryInput struct {
	// AssignmentID uuid.UUID `json:"assignmentId" binding:"required"`
	Coordinates LocationCoordinatesInput `json:"coords" json:""`
	Note        string                   `json:"note" binding:""`
	Description string                   `json:"description" binding:""`
	Active      bool                     `json:"active" binding:""`
}

type RespondentSessionLocationEntryCreationInput struct {
	RespondentSessionLocationEntryInput
}

type RespondentSessionLocationEntryUpdateInput struct {
	Note        *string `json:"note" binding:""`
	Description *string `json:"description" binding:""`
	Active      *bool   `json:"active" binding:""`
}
