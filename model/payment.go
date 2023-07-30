package model

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	// Code     string    `gorm:"" json:"code"`
	CurrencyCode string `gorm:"type:char(3)" json:"currencyCode"`
	Amount       int64  `json:"total,omitempty"`

	//THE ORDER ENTITY
	OrderID *uuid.UUID `gorm:"not null" json:"orderId,omitempty"`
	Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"order,omitempty"`

	Status string `json:"status,omitempty"`

	//Payment Method
	PaymentMethodID *uuid.UUID     `gorm:"" json:"paymentMethodId,omitempty"`
	PaymentMethod   *PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"paymentMethod,omitempty"`

	//STORE PAYMENT DETAILS LIKE CREDIT CARD INFOR, PAPAL_ID in a map
	Details map[string]interface{} `gorm:"" json:"details"`

	//TIMESTAMPs
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}
