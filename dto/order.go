package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
)

type OrderCreationInput struct {
	UserID          *uuid.UUID     `json:"userId" binding:""`
	PaymentMethodID *uuid.UUID     `json:"paymentMethodId" binding:""`
	Items           []RequestInput `json:"requests" binding:""`
	Situations      []uuid.UUID    `json:"situations" binding:"required"`
	Description     *string        `json:"description" binding:""`
}

type OrderUpdateInput struct {
	PaymentMethodID *uuid.UUID `json:"paymentMethodId" binding:""`
}

type RespondentOrderCompletionInput struct {
	ClientPaidCash bool `json:"clientPaidCash" binding:""`
}

type OrderOutput struct {
	ID               uuid.UUID `json:"id,omitempty"`
	Code             string    `json:"code"`
	ItemsTotal       uint64    `json:"itemsTotal,omitempty"`
	AdjustmentsTotal int64     `json:"adjustmentTotal,omitempty"`
	Total            uint64    `json:"total,omitempty"`

	Items       *[]model.Request        `json:"items,omitempty"`
	Adjustments []model.OrderAdjustment `json:"adjustments,omitempty"`

	//
	Status model.OrderStatus `json:"status,omitempty"`

	//ASSOCIATED USER ACCOUNT
	UserID *uuid.UUID  `json:"userId,omitempty"`
	User   *model.User `json:"user,omitempty"`

	//REQUESTING USER ACCOUNT
	FulfilmentID *uuid.UUID             `json:"fulfilmentId,omitempty"`
	Fulfilment   *model.OrderFulfilment `json:"fulfilment,omitempty"`

	//Payment Method
	PaymentMethodID *uuid.UUID           `json:"paymentMethodId,omitempty"`
	PaymentMethod   *model.PaymentMethod `json:"paymentMethod,omitempty"`
	ClientPaidCash  *bool                `gorm:"" json:"clientPaidCash"`

	//Payment Method
	PaymentID *uuid.UUID          `json:"paymentId,omitempty"`
	Payment   *model.OrderPayment ` json:"payment,omitempty"` //

	// SITUATIONS
	Situations  *[]*model.MotoristRequestSituation `json:"situations,omitempty"` //`json:"situations,omitempty"`
	Description *string                            `json:"description,omitempty"`
	// OrderMotoristRequestSituation []*model.OrderMotoristRequestSituation ``
	//TIMESTAMPS
	CheckoutCompletedAt *time.Time `json:"checkoutCompletedAt,omitempty"`
	CreatedAt           *time.Time `json:"createdAt,omitempty"`
	UpdatedAt           *time.Time `json:"updatedAt,omitempty"`
}

func CreateOrderOutput(order model.Order) (output OrderOutput) {

	var Situations []*model.MotoristRequestSituation // = make([]*model.MotoristRequestSituation, 1)
	if order.OrderMotoristRequestSituations != nil {
		for _, sit := range order.OrderMotoristRequestSituations {
			Situations = append(Situations, sit.Situation)
		}
	}

	return OrderOutput{
		ID:   order.ID,
		Code: order.Code,

		ItemsTotal:       order.ItemsTotal,
		AdjustmentsTotal: order.AdjustmentsTotal,
		Total:            order.Total,

		Items:       order.Items,
		Adjustments: order.Adjustments,

		//
		Status: order.Status,

		//ASSOCIATED USER ACCOUNT
		UserID: order.UserID,
		User:   order.User,

		//REQUESTING USER ACCOUNT
		FulfilmentID: order.FulfilmentID,
		Fulfilment:   order.Fulfilment,

		//Payment Method
		PaymentMethodID: order.PaymentMethodID,
		PaymentMethod:   order.PaymentMethod,

		//Payment Method
		PaymentID:      order.PaymentID,
		Payment:        order.Payment,
		ClientPaidCash: order.ClientPaidCash,

		// SITUATIONS
		Situations:  &Situations, //`json:"situations,omitempty"`
		Description: order.Description,
		// OrderMotoristRequestSituation []*model.OrderMotoristRequestSituation ``
		//TIMESTAMPS
		CheckoutCompletedAt: order.CheckoutCompletedAt,
		CreatedAt:           order.CreatedAt,
		UpdatedAt:           order.UpdatedAt,
	}
}
