package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/initializers"
	"gorm.io/gorm"
)

type RespondentAccessProductPrice struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`

	Label       *string `gorm:"type:varchar(128);" json:"label,omitempty"`
	Description *string `gorm:"type:varchar(225);" json:"description,omitempty"`
	//
	UnitPrice *uint64 `gorm:"not null" json:"unitPrice"`
	Upto      *uint64 `gorm:"bigint(24);not null;" json:"upTo"` //uniqueIndex:UNIQUE_PRODUCT_PRICE_UPTO_AND_TYPE

	ProductType *string `gorm:"varchar(255);" json:"productType"` // PURCHASE  or SUBSCRIPTION//uniqueIndex:UNIQUE_PRODUCT_PRICE_UPTO_AND_TYPE
	//
	StripePriceID *string `gorm:"not null" json:"stripePriceId,omitempty"`
	Active        *bool   `gorm:"not null;default:true" json:"active,omitempty"`

	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

// func (price *RespondentAccessProductPrice) Increment(quantity int64) {
// 	if price.Price == nil {
// 		var Val int64 = 0
// 		price.Price = &Val
// 	}
// 	*price.Price += quantity
// }

func (price *RespondentAccessProductPrice) BeforeCreate(tx *gorm.DB) (err error) {
	price.ID = uuid.New()
	price.CreatedAt = time.Now()
	price.UpdatedAt = time.Now()

	return nil
}

func (price *RespondentAccessProductPrice) ProductID() (productId *string) {
	if *price.ProductType == ACCESS_PRODUCT_TYPE_PURCHASE {
		productId = &initializers.CONFIG.STRIPE_RESPONDENT_PURCHASE_PRODUCT_ID
	} else if *price.ProductType == ACCESS_PRODUCT_TYPE_SUBSCIPTION {
		productId = &initializers.CONFIG.STRIPE_RESPONDENT_SUBSCRIPTION_PRODUCT_ID
	} else {
		productId = nil
	}
	return productId
}

type AccessProductType = string

const ACCESS_PRODUCT_TYPE_PURCHASE AccessProductType = "PURCHASE"
const ACCESS_PRODUCT_TYPE_SUBSCIPTION AccessProductType = "SUBSCRIPTION"
