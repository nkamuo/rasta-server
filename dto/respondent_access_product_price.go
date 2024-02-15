package dto

// import "github.com/google/uuid"

type RespondentPurchaseSessionCheckoutInput struct {
	PriceID  string `json:"priceId" binding:"required"`
	Quantity *int64 `json:"quantity" binding:""`
}
