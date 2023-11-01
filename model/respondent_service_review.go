package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RespondentServiceReview struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	//
	Rating      uint8   `gorm:"" json:"rating"`
	Description *string `gorm:"" json:"description,omitempty"`
	//
	RespondentID *uuid.UUID  `gorm:"not null" json:"respondentId,omitempty"`
	Respondent   *Respondent `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"`
	//
	OrderID *uuid.UUID `gorm:"not null" json:"orderId,omitempty"`
	Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order,omitempty"`

	ArrivedOnTime *bool `gorm:"" json:"arrivedOnTime"`
	//
	Published *bool     `gorm:"default:false;not null" json:"published"`
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (review *RespondentServiceReview) BeforeCreate(tx *gorm.DB) (err error) {
	review.ID = uuid.New()

	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()
	return nil
}
