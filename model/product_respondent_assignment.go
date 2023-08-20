package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRespondentAssignment struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`

	//ASSOCIATED USER ACCOUNT
	RespondentID *uuid.UUID  `gorm:"not null" json:"respondentId,omitempty"`
	Respondent   *Respondent `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"` // BelongsToMany Association - These mappings may be optional
	//Product OF THE REPSONDANT
	ProductID *uuid.UUID `gorm:"not null" json:"productId,omitempty"`
	Product   *Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product,omitempty"` // BelongsToMany Association - These mappings may be optional
	//TOGGLES
	Active                  *bool `gorm:"default:false;not null" json:"active"`
	AllowRespondentActivate *bool `gorm:"default:false;not null" json:"allowRespondentActivate"` //Allows the respondent to togle the Active state
	//DOCUMENTATION
	Note        string `gorm:"LENGTH(255);" json:"note,omitempty"`
	Description string `gorm:"not null" json:"description,omitempty"`
	// TIMESTAMPS
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (assignment *ProductRespondentAssignment) BeforeCreate(tx *gorm.DB) (err error) {
	assignment.ID = uuid.New()
	return nil
}
