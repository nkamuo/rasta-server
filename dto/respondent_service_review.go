package dto

import "github.com/google/uuid"

type RespondentServiceReviewCreationInput struct {
	Rating uint8 `json:"rating" binding:"required"`
	//
	ArrivedOnTime *bool `json:"arrivedOnTime,omitempty" binding:""`
	//
	RequestID uuid.UUID  `json:"requestId,omitempty" binding:"required"`
	AuthorID  *uuid.UUID `gorm:"not null" json:"author,omitempty" binding:""` //THE ID of the user who made the review
	//
	Description *string `json:"description,omitempty" binding:""`
	Published   *bool   `json:"published" binding:""`
}

type RespondentServiceReviewUpdateInput struct {
	Rating      *uint8  `json:"rating" binding:"required"`
	Description *string `json:"description,omitempty" binding:""`
	Published   *bool   `json:"published" binding:""`
}
