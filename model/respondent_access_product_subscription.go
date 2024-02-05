package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RespondentAccessProductSubscription struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`

	//ASSOCIATED USER ACCOUNT
	RespondentID *uuid.UUID  `gorm:"not null" json:"respondentId,omitempty"`
	Respondent   *Respondent `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"` // BelongsToMany Association - These mappings may be optional
	//TOGGLES
	Active *bool `gorm:"default:false;not null" json:"active"`
	//DOCUMENTATION
	StripeProductID      *string `gorm:"LENGTH(255);" json:"stripeProductId,omitempty"`
	StripeProductPriceID *string `gorm:"LENGTH(255);" json:"stripeProductPriceId,omitempty"`
	StripeSubscriptionID *string `gorm:"LENGTH(255);" json:"stripeSubscriptionId,omitempty"`
	//
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (subscription *RespondentAccessProductSubscription) BeforeCreate(tx *gorm.DB) (err error) {
	subscription.ID = uuid.New()
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()

	return nil
}
