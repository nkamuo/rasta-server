package dto

import "github.com/google/uuid"

type CompanyCreationInput struct {
	OperatorUserID uuid.UUID `json:"operatorUserId" binding:"required"`
	Title          string    `json:"title" binding:"required"`
	LicenseNumber  string    `json:"licenseNumber" binding:"required"`
	Description    string    `json:"description" binding:"required"`
	Category       string    `json:"category" binding:"required"`
	Active         bool      `json:"active" binding:""`
	Published      bool      `json:"published" binding:""`
}

type CompanyUpdateInput struct {
	OperatorUserID *uuid.UUID `json:"operatorUserId" binding:""`
	LicenseNumber  *string    `json:"licenseNumber" binding:""`
	Title          *string    `json:"title" binding:""`
	Description    *string    `json:"description" binding:""`
	Category       *string    `json:"category" binding:""`
	Active         *bool      `json:"active" binding:""`
	Published      *bool      `json:"published" binding:""`
}
