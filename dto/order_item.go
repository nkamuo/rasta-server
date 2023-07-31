package dto

import "github.com/google/uuid"

type OrderItemInput struct {
	ProductID   *uuid.UUID `gorm:"" json:"productId,omitempty"`
	UnitPrice   *int64     `json:"unitPrice"`
	Quantity    *int64     `json:"quantity,omitempty"`
	Origin      *string    `json:"originId"`
	Destination *string    `json:"destinationId"`
	Note        *string    `json:"note"`
}

// type OrderItemUpdateInput struct {
// 	PaymentMethodID *uuid.UUID `json:"paymentMethodId" binding:""`
// }
