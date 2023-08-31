package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RespondentServiceReview struct {
	ID     uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Rating uint8     `gorm:"" json:"rating"`
	//
	RespondentID *uuid.UUID `gorm:"not null" json:"respondentId,omitempty"`
	Respondent   *FuelType  `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"`
	//
	RequestID *uuid.UUID `gorm:"not null" json:"requestId,omitempty"`
	Request   *Request   `gorm:"foreignKey:RequestID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"request,omitempty"`
	//
	AuthorID *uuid.UUID `gorm:"not null" json:"authorId,omitempty"`
	Author   *User      `gorm:"foreignKey:AuthorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"author,omitempty"`
	//
	Description *string   `gorm:"" json:"description,omitempty"`
	Published   *bool     `gorm:"default:false;not null" json:"published"`
	CreatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (review *RespondentServiceReview) BeforeCreate(tx *gorm.DB) (err error) {
	review.ID = uuid.New()

	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()
	return nil
}
