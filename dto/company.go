package dto

import "github.com/google/uuid"

type CompanyCreationInput struct {
	OperatorUserID uuid.UUID `json:"operatorUserId" binding:"required"`
	Title          string    `json:"title" binding:"required"`
	Description    string    `json:"description" binding:"required"`
	Category       string    `json:"category" binding:"required"`
}

type CompanyUpdateInput struct {
	Title       string `json:"title" binding:"optional"`
	Description string `json:"description" binding:"optional"`
	Category    string `json:"category" binding:"optional"`
}
