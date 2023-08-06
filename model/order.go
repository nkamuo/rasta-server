package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID               uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Code             string    `gorm:"" json:"code"`
	ItemsTotal       uint64    `json:"itemsTotal,omitempty"`
	AdjustmentsTotal uint64    `json:"adjustmentTotal,omitempty"`
	Total            uint64    `json:"total,omitempty"`

	Items *[]OrderItem `gorm:"foreignKey:OrderID"`

	Status string `json:"status,omitempty"`

	//ASSOCIATED USER ACCOUNT
	UserID *uuid.UUID `gorm:"" json:"userId,omitempty"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`

	//Payment Method
	PaymentMethodID *uuid.UUID     `gorm:"" json:"paymentMethodId,omitempty"`
	PaymentMethod   *PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"paymentMethod,omitempty"`

	//TIMESTAMPs
	CheckoutCompletedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"checkoutCompletedAt,omitempty"`
	CreatedAt           time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt           time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) (err error) {
	order.ID = uuid.New()

	var total uint64
	for _, item := range *order.Items {
		total += (item.Quantity * item.Rate)
	}
	order.Total = total

	return nil
}