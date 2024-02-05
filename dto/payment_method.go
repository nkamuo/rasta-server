package dto

import (
	"github.com/nkamuo/rasta-server/model"
)

type PaymentMethodCreationInput struct {
	// UserID      uuid.UUID                   `json:"userId" binding:""`
	// Category    model.PaymentMethodCategory `json:"category" binding:"required"`
	// Details     model.JSON/*map[string]interface{}*/ `json:"details" binding:""`
	// Description string `json:"description" binding:""`
	// Active      bool   `json:"active" binding:""`
	Token string `json:"token" binding:"required"`
}

type PaymentMethodUpdateInput struct {
	Category    *model.PaymentMethodCategory `json:"category" binding:"required"`
	Details     *model.JSON/**map[string]interface{}*/ `json:"details" binding:""`
	Description *string `json:"description" binding:""`
	Active      *bool   `json:"active" binding:""`
}

type SelectDefaultPaymentMethodInput struct {
	PaymentMethodID string `json:"paymentMethodId" binding:"required"`
}
