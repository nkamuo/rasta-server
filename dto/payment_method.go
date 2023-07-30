package dto

import (
	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
)

type PaymentMethodCreationInput struct {
	UserID      uuid.UUID                   `json:"userId" binding:""`
	Category    model.PaymentMethodCategory `json:"category" binding:"required"`
	Details     map[string]interface{}      `json:"details" binding:""`
	Description string                      `json:"description" binding:""`
	Active      bool                        `json:"active" binding:""`
}

type PaymentMethodUpdateInput struct {
	Category    *model.PaymentMethodCategory `json:"category" binding:"required"`
	Details     *map[string]interface{}      `json:"details" binding:""`
	Description *string                      `json:"description" binding:""`
	Active      *bool                        `json:"active" binding:""`
}
