package dto

import "github.com/google/uuid"

type UserRegistrationInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	//PROFILE INFORMATION
	FirstName    string  `json:"firstName" binding:"required"`
	LastName     string  `json:"lastName" binding:"required"`
	Phone        string  `json:"phone" binding:"required"`
	ReferrerCode *string `json:"referrerCode" binding:""`
	//RESPONDENT INFORMATION
	IsRespondent      bool       `json:"isRespondent" binding:""`
	RespondentPlaceId *uuid.UUID `json:"respondentPlaceId" binding:""`
}

type UserFormLoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
