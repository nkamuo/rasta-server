package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RespondentSession struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`

	//ASSOCIATED USER ACCOUNT
	RespondentID *uuid.UUID  `gorm:"not null" json:"respondentId,omitempty"`
	Respondent   *Respondent `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"` // BelongsToMany Association - These mappings may be optional
	//TOGGLES
	Active *bool `gorm:"default:false;not null" json:"active"`
	//DOCUMENTATION
	Note        string `gorm:"LENGTH(255);" json:"note,omitempty"`
	Description string `gorm:"not null" json:"description,omitempty"`
	//
	StartingCoordinates LocationCoordinates                `gorm:"embedded;columnPrefix:starting_coords_" json:"startingCoords"`
	Assignments         []RespondentSessionAssignedProduct `gorm:"foreignKey:SessionID" json:"assignments"`
	// Assignments         []ProductRespondentAssignment `gorm:"many2many:respondent_session_active_product_assignments;"`
	// TIMESTAMPS
	StartedAt *time.Time `gorm:"" json:"startedAt,omitempty"`
	EndedAt   *time.Time `gorm:"" json:"endedAt,omitempty"`
	//
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (assignment *RespondentSession) BeforeCreate(tx *gorm.DB) (err error) {
	assignment.ID = uuid.New()
	return nil
}
