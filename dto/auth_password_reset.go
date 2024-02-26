package dto

type UserFormResetPasswordCodeRequestInput struct {
	Email string `json:"email" binding:"required,email"`
}

type UserFormResetPasswordCodeValidationInput struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type UserFormResetPasswordInput struct {
	Code       string `json:"code" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	ResetToken string `json:"resetToken" binding:"required"`
	Password   string `json:"password" binding:"required"`
}
