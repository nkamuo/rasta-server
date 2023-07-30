package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PaymentMethod struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	//ASSOCIATED USER ACCOUNT
	UserID *uuid.UUID `gorm:"unique" json:"userId,omitempty"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`

	Category PaymentMethodCategory `gorm:"varchar(100)" json:"category,omitempty"`

	//STORE PAYMENT DETAILS LIKE CREDIT CARD INFOR, PAPAL_ID in a map
	Details map[string]interface{} `gorm:"" json:"details"`

	Description string `gorm:"" json:"description,omitempty"`
	Active      bool   `gorm:"default:false;not null" json:"active"`

	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

type PaymentMethodCategory = string

const (
	PAYMENT_METHOD_CREDIT_CARD ProductCategory = "CREDIT_CARD"
	PAYMENT_METHOD_CASH        ProductCategory = "CASH"
)

func ValidatePaymentMethodCategory(category PaymentMethodCategory) (err error) {
	switch category {
	case PAYMENT_METHOD_CREDIT_CARD:
		return nil
	case PAYMENT_METHOD_CASH:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported payment method type \"%s\"", category))
}
