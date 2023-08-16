package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID               uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Code             string    `gorm:"" json:"code"`
	ItemsTotal       uint64    `json:"itemsTotal,omitempty"`
	AdjustmentsTotal int64     `json:"adjustmentTotal,omitempty"`
	Total            uint64    `json:"total,omitempty"`

	Items       *[]Request        `gorm:"foreignKey:OrderID"`
	Adjustments []OrderAdjustment `gorm:"foreignKey:OrderID"`

	Status string `json:"status,omitempty"`

	//ASSOCIATED USER ACCOUNT
	UserID *uuid.UUID `gorm:"" json:"userId,omitempty"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`

	//Payment Method
	PaymentMethodID *uuid.UUID     `gorm:"" json:"paymentMethodId,omitempty"`
	PaymentMethod   *PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"paymentMethod,omitempty"`

	//Payment Method
	PaymentID *uuid.UUID    `gorm:"" json:"paymentId,omitempty"`
	Payment   *OrderPayment `gorm:"foreignKey:PaymentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"payment,omitempty"`

	//TIMESTAMPs
	CheckoutCompletedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"checkoutCompletedAt,omitempty"`
	CreatedAt           time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt           time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) (err error) {
	order.ID = uuid.New()
	itemTotal := order.CalculateItemTotal()
	adjustmentTotal := order.CalculateAdjustmentTotal()
	if adjustmentTotal < 0 {
		if uint64(adjustmentTotal) > itemTotal {
			return errors.New(fmt.Sprintf("Negetive Adjustment total of [%d] is greater than order total [%d]", adjustmentTotal, itemTotal))
		}
	}
	order.ItemsTotal = itemTotal
	order.AdjustmentsTotal = adjustmentTotal
	order.Total = itemTotal + uint64(adjustmentTotal)
	return nil
}

func (order *Order) CalculateItemTotal() (itemTotal uint64) {
	var total uint64
	for _, item := range *order.Items {
		total += (item.Quantity * item.Rate)
	}
	return total
}

func (order *Order) CalculateTotal() (total uint64, err error) {
	order.ID = uuid.New()
	itemTotal := order.CalculateItemTotal()
	adjustmentTotal := order.CalculateAdjustmentTotal()
	if adjustmentTotal < 0 {
		if uint64(adjustmentTotal) > itemTotal {
			return 0, errors.New(fmt.Sprintf("Negetive Adjustment total of [%d] is greater than order total [%d]", adjustmentTotal, itemTotal))
		}
	}
	total = itemTotal + uint64(adjustmentTotal)
	return total, nil
}

func (order *Order) CalculateAdjustmentTotal() (adjustmentTotal int64) {
	var total int64
	for _, adjustment := range order.Adjustments {
		total += adjustment.Amount
	}
	return total
}

func (order *Order) AddAdjustment(adjustment OrderAdjustment) (err error) {
	adjustment.OrderID = &order.ID
	order.Adjustments = append(order.Adjustments, adjustment)
	return nil
}

const SERVICE_FEE_ADJUSTMENT_CODE = "SERVICE_FEE"
