package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RespondentSessionAssignedProduct struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`

	//ASSOCIATED USER ACCOUNT
	SessionID *uuid.UUID         `gorm:"not null" json:"sessionId,omitempty"`
	Session   *RespondentSession `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"session,omitempty"` // BelongsToMany Association - These mappings may be optional
	//Product OF THE REPSONDANT
	AssignmentID *uuid.UUID                   `gorm:"not null" json:"assignmentId,omitempty"`
	Assignment   *ProductRespondentAssignment `gorm:"foreignKey:AssignmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assignment,omitempty"` // BelongsToMany Association - These mappings may be optional
	//TOGGLES
	Active *bool `gorm:"default:false;not null" json:"active"`
	//DOCUMENTATION
	Note        string `gorm:"LENGTH(255);" json:"note,omitempty"`
	Description string `gorm:"not null" json:"description,omitempty"`
	//
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (assignment *RespondentSessionAssignedProduct) BeforeCreate(tx *gorm.DB) (err error) {
	assignment.ID = uuid.New()

	assignment.CreatedAt = time.Now()
	assignment.UpdatedAt = time.Now()
	return nil
}
