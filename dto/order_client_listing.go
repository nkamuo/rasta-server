package dto

type ClientOrderConfirmationRequest struct {
	Rating      uint8   `json:"rating" binding:"required,lte=5"`
	Description *string `json:"description" binding:""`
}
