package dto

import "github.com/google/uuid"

type OrderCreationInput struct {
	UserID          *uuid.UUID     `json:"userId" binding:""`
	PaymentMethodID *uuid.UUID     `json:"paymentMethodId" binding:""`
	Items           []RequestInput `json:"requests" binding:""`
	Situations      []uuid.UUID    `json:"situations" binding:"required"`
}

type OrderUpdateInput struct {
	PaymentMethodID *uuid.UUID `json:"paymentMethodId" binding:""`
}
