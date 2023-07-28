package dto

import (
	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
)

type ProductCreationInput struct {
	PlaceID     uuid.UUID             `json:"placeId" binding:"required"`
	Category    model.ProductCategory `json:"category" binding:"required"`
	Label       string                `json:"label" binding:"required"`
	IconImage   string                `json:"iconImage" binding:"required"`
	CoverImage  string                `json:"coverImage" binding:""`
	Title       string                `json:"title" binding:""`
	Description string                `json:"description" binding:""`
}

type ProductUpdateInput struct {
	// PlaceID     uuid.UUID             `json:"placeId" binding:"required"`
	// Category    model.ProductCategory `json:"category" binding:"required"`
	Label       string `json:"label" binding:""`
	IconImage   string `json:"iconImage" binding:""`
	CoverImage  string `json:"coverImage" binding:""`
	Title       string `json:"title" binding:""`
	Description string `json:"description" binding:""`
}
