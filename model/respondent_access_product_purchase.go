package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RespondentAccessProductPurchase struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`

	//ASSOCIATED USER ACCOUNT
	RespondentID *uuid.UUID  `gorm:"not null" json:"respondentId,omitempty"`
	Respondent   *Respondent `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"` // BelongsToMany Association - These mappings may be optional
	//Purchase
	// Purchase *int64 `gorm:"default:0;not null" json:"balance"`
	// The total number of request the user has paid for and can handle
	//TODO: Decremete this number each time he handles a request

	StripeCheckoutSessionID *string                       `gorm:"not null" json:"stripeCheckoutSessionId,omitempty"`
	StripePriceID           *string                       `gorm:"not null" json:"stripePriceId,omitempty"`
	PriceID                 *uuid.UUID                    `gorm:"type:char(36);" json:"priceId,omitempty"`
	Price                   *RespondentAccessProductPrice `gorm:"foreignKey:PriceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"price,omitempty"`
	// The price ID of the product purchased
	Quantity  *int64 `gorm:"default:1;not null" json:"quantity,omitempty"`
	Succeeded *bool  `gorm:"default:false;not null" json:"succeeded,omitempty"`
	Cancelled *bool  `gorm:"default:false;not null" json:"cancelled,omitempty"`

	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (balance *RespondentAccessProductPurchase) BeforeCreate(tx *gorm.DB) (err error) {
	balance.ID = uuid.New()
	balance.CreatedAt = time.Now()
	balance.UpdatedAt = time.Now()

	return nil
}
