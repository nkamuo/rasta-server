package dto

// import "github.com/google/uuid"

type RespondentPurchaseAdminInput struct {
	PriceID  string `json:"priceId" binding:"required"`
	Quantity *int64 `json:"quantity" binding:""`
	Commit   *bool  `json:"commit" binding:""`
}
