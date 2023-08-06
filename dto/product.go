package dto

import (
	"github.com/google/uuid"
	"github.com/nkamuo/rasta-server/model"
)

type ProductCreationInput struct {
	PlaceID     uuid.UUID             `json:"placeId" binding:"required"`
	Category    model.ProductCategory `json:"category" binding:"required"`
	Bundled     bool                  `json:"bundled" binding:""`
	Published   bool                  `json:"published" binding:""`
	Rate        uint64                `json:"rate" binding:""`
	Label       string                `json:"label" binding:"required"`
	Title       string                `json:"title" binding:""`
	IconImage   string                `json:"iconImage" binding:""`
	CoverImage  string                `json:"coverImage" binding:""`
	Description string                `json:"description" binding:""`
}

type ProductUpdateInput struct {
	// PlaceID     uuid.UUID             `json:"placeId" binding:"required"`
	// Category    model.ProductCategory `json:"category" binding:"required"`
	Published   *bool   `json:"published" binding:""`
	Bundled     *bool   `json:"bundled" binding:""`
	Rate        *uint64 `json:"rate" binding:""`
	Label       *string `json:"label" binding:""`
	Title       *string `json:"title" binding:""`
	IconImage   *string `json:"iconImage" binding:""`
	CoverImage  *string `json:"coverImage" binding:""`
	Description *string `json:"description" binding:""`
}
