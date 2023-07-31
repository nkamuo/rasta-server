package dto

import "github.com/google/uuid"

type OrderCreationInput struct {
	UserID          *uuid.UUID       `json:"userId" binding:""`
	PaymentMethodID *uuid.UUID       `json:"paymentMethodId" binding:""`
	Items           []OrderItemInput `json:"items" binding:""`
}

type OrderUpdateInput struct {
	PaymentMethodID *uuid.UUID `json:"paymentMethodId" binding:""`
}
