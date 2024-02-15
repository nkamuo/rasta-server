package dto

// import "github.com/google/uuid"

type RespondentAccessProductPriceCreationInput struct {
	// StripePriceID string  `json:"stripePriceID" binding:"required"`
	UnitPrice   uint64  `json:"unitPrice" binding:"required"`
	Mode        string  `json:"mode" binding:"required"`
	UpTo        *uint64 `json:"upTo" binding:"required"`
	Active      *bool   `json:"active" binding:""`
	Label       *string `json:"label" binding:""`
	Description *string `json:"description" binding:""`
	// Quantity *int64 `json:"quantity" binding:""`
}

type RespondentAccessProductPriceUpdateInput struct {
	// StripePriceID *string `json:"stripePriceID" binding:""`
	UnitPrice   *uint64 `json:"unitPrice" binding:""`
	Mode        *string `json:"mode" binding:""`
	UpTo        *uint64 `json:"upTo" binding:""`
	Label       *string `json:"label" binding:""`
	Active      *bool   `json:"active" binding:""`
	Description *string `json:"description" binding:""`
	// Quantity *int64 `json:"quantity" binding:""`
}
