package dto

type UserPasswordResetRequestInput struct {
	Email string `json:"email" binding:"required"`
}
