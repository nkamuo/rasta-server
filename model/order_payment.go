package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderPayment struct {
	ID     uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Status string    `gorm:"type:varchar(32);not null;" json:"status,omitempty"`
	Amount int64     `json:"amount,omitempty"`
	// Code        string    `gorm:"not null;varchar(32)" json:"code,omitempty"`
	Title       string  `gorm:"varchar(64)" json:"title,omitempty"`
	Description *string `gorm:"varchar(225)" json:"description,omitempty"`

	StripeID           *string `gorm:"varchar(225)" json:"stripeId,omitempty"`
	StripeClientSecret *string `gorm:"varchar(225)" json:"clientSecret,omitempty"`

	//THE MAIN ORDER ENTITY
	OrderID *uuid.UUID `gorm:"not null" json:"orderId,omitempty"`
	Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order,omitempty"`
	//
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (payment *OrderPayment) BeforeCreate(tx *gorm.DB) (err error) {
	payment.ID = uuid.New()
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	return nil
}
