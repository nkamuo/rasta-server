package dto

import "github.com/google/uuid"

type CreateCustomerInput struct {
	UserId uuid.UUID `json:"userId" binding:"required"`
	// Title       string `json:"title" binding:"required"`
	// Description string `json:"description" binding:"required"`
}
