package model

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	ORDER_STATUS_DRAFT                OrderStatus = "draft"
	ORDER_STATUS_PENDING              OrderStatus = "pending"
	ORDER_STATUS_RESPONDENT_ASSIGNED  OrderStatus = "assigned"
	ORDER_STATUS_RESPONDENT_ARRIVED   OrderStatus = "arrived"
	ORDER_STATUS_CLIENT_CONFIRMED     OrderStatus = "client_confirmed"
	ORDER_STATUS_CLIENT_REJECTED      OrderStatus = "client_rejected"
	ORDER_STATUS_RESPONDENT_CONFIRMED OrderStatus = "responder_confirmed"
	ORDER_STATUS_RESPONDENT_REJECTED  OrderStatus = "responder_rejected"
	ORDER_STATUS_CANCELLED            OrderStatus = "cancelled"
	ORDER_STATUS_COMPLETED            OrderStatus = "completed"
	//
	ORDER_STATUS_COMPLETED_BY_RESPONDENT OrderStatus = "responder_done"
)

type Order struct {
	ID               uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Code             string    `gorm:"" json:"code"`
	ItemsTotal       uint64    `json:"itemsTotal,omitempty"`
	AdjustmentsTotal int64     `json:"adjustmentTotal,omitempty"`
	Total            uint64    `json:"total,omitempty"`

	Items       *[]Request        `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Adjustments []OrderAdjustment `gorm:"foreignKey:OrderID" json:"adjustments,omitempty"`

	//
	Status OrderStatus `gorm:"type:varchar(32);not null;default:'draft'" json:"status,omitempty"`

	//ASSOCIATED USER ACCOUNT
	UserID *uuid.UUID `gorm:"" json:"userId,omitempty"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`

	//REQUESTING USER ACCOUNT
	FulfilmentID *uuid.UUID       `gorm:"unique" json:"fulfilmentId,omitempty"`
	Fulfilment   *OrderFulfilment `gorm:"foreignKey:FulfilmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"fulfilment,omitempty"`

	//Payment Method
	PaymentMethodID *uuid.UUID     `gorm:"" json:"paymentMethodId,omitempty"`
	PaymentMethod   *PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"paymentMethod,omitempty"`

	//Payment Method
	PaymentID *uuid.UUID    `gorm:"" json:"paymentId,omitempty"`
	Payment   *OrderPayment ` json:"payment,omitempty"` //gorm:"foreignKey:PaymentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"
	//
	ClientPaidCash *bool `gorm:"" json:"clientPaidCash"`

	// SITUATIONS
	// Situations                    *[]*MotoristRequestSituation     `gorm:""` //`gorm:"many2many:order_motorist_situations;ForeignKey:ID;References:ID;joinForeignKey:order_id;joinReferences:motorist_request_situation_id;" json:"situations,omitempty"`
	OrderMotoristRequestSituations []*OrderMotoristRequestSituation `gorm:""`
	Description                    *string                          `gorm:"" json:"description,omitempty"`
	//
	Review *RespondentServiceReview `gorm:"" json:"review,omitempty"`
	//TIMESTAMPS
	CheckoutCompletedAt *time.Time `gorm:";" json:"checkoutCompletedAt,omitempty"`
	CreatedAt           *time.Time `gorm:"not null;" json:"createdAt,omitempty"`
	UpdatedAt           *time.Time `gorm:"ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

type OrderMotoristRequestSituation struct {
	OrderID *uuid.UUID `gorm:"primaryKey" json:"orderId,omitempty"`
	Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order,omitempty"`
	//

	SituationID *uuid.UUID                `gorm:"primaryKey" json:"situationId,omitempty"`
	Situation   *MotoristRequestSituation `gorm:"foreignKey:SituationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"situation,omitempty"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	order.ID = uuid.New()
	order.Code = generateOrderCode()
	order.Status = ORDER_STATUS_DRAFT
	order.CreatedAt = &now
	// order.UpdatedAt = &now
	//
	itemTotal := order.CalculateItemTotal()
	adjustmentTotal := order.CalculateAdjustmentTotal()
	total, err := order.CalculateTotal()
	if nil != err {
		return err
	}

	order.ItemsTotal = itemTotal
	order.AdjustmentsTotal = adjustmentTotal
	order.Total = total
	return nil
}

func (order *Order) CalculateItemTotal() (itemTotal uint64) {
	items := order.Items
	if nil == items {
		items = &[]Request{}
	}
	for _, item := range *items {
		quantity := item.Quantity
		if quantity == 0 {
			quantity = 1
		}
		itemTotal += (quantity * item.Rate)
	}
	return itemTotal
}

func (order *Order) GetTotal() (total uint64, err error) {
	if order.Total != 0 {
		return order.Total, nil
	}
	return order.CalculateTotal()
}

func (order *Order) CalculateTotal() (total uint64, err error) {
	itemTotal := order.CalculateItemTotal()
	adjustmentTotal := order.CalculateAdjustmentTotal()
	if adjustmentTotal < 0 {
		if uint64(adjustmentTotal) > itemTotal {
			return 0, errors.New(fmt.Sprintf("Negetive Adjustment total of [%d] is greater than order total [%d]", adjustmentTotal, itemTotal))
		}
		total = itemTotal - uint64(adjustmentTotal)
	} else {
		total = itemTotal + uint64(adjustmentTotal)
	}
	return total, nil
}

func (order *Order) CalculateAdjustmentTotal() (adjustmentTotal int64) {
	for _, adjustment := range order.Adjustments {
		adjustmentTotal += adjustment.Amount
	}
	return adjustmentTotal
}

func (order *Order) AddAdjustment(adjustment OrderAdjustment) (err error) {
	adjustment.OrderID = &order.ID
	order.Adjustments = append(order.Adjustments, adjustment)
	return nil
}

func (order *Order) GetPrimaryLocation() (location *Location, err error) {

	if len(*order.Items) < 1 {
		return nil, errors.New("Order must have at least one Item")
	}
	item := (*order.Items)[0]

	if item.Origin != nil {
		return item.Origin, nil
	}

	if item.Destination != nil {
		return item.Destination, nil
	}

	return nil, errors.New("Order Item must have destination")

}

const SERVICE_FEE_ADJUSTMENT_CODE = "SERVICE_FEE"

func generateOrderCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	// rand.NewRand()

	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
