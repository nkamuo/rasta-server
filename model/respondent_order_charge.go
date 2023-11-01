package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RespondentOrderCharge struct {
	ID          uuid.UUID          `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Amount      uint64             `gorm:"" json:"amount"`
	Label       string             `gorm:"type:varchar(128);not null" json:"label,omitempty"`
	Description string             `gorm:"not null" json:"description,omitempty"`
	Status      OrderEarningStatus `gorm:"type:varchar(16);not null" json:"status,omitempty"`

	// STRIPE ID
	StripePaymentID           *string `gorm:"varchar(225)" json:"stripePaymentId,omitempty"`
	StripePaymentClientSecret *string `gorm:"varchar(225)" json:"stripPaymentClientSecret,omitempty"`
	//THE RESPONDENT EXECUTING THIS TASK
	RequestID *uuid.UUID `gorm:"not null" json:"requestId,omitempty"`
	Request   *Request   `gorm:"foreignKey:RequestID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"request,omitempty"`

	CommittedAt *time.Time `gorm:"" json:"committedAt,omitempty"`
	//TIMESTAMPs
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
	//THE RESPONDENT EXECUTING THIS TASK
	RespondentID *uuid.UUID  `gorm:"not null" json:"respondentId,omitempty"`
	Respondent   *Respondent `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"`
}

func (charge *RespondentOrderCharge) BeforeCreate(tx *gorm.DB) (err error) {
	charge.ID = uuid.New()
	charge.CreatedAt = time.Now()
	charge.UpdatedAt = time.Now()
	return nil
}
