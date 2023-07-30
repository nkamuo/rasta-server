package dto

import (
	"github.com/nkamuo/rasta-server/model"
)

type VehicleModelCreationInput struct {
	Category    model.VehicleCategory `json:"category" binding:"required"`
	Label       string                `json:"label" binding:"required"`
	IconImage   string                `json:"iconImage" binding:"required"`
	CoverImage  string                `json:"coverImage" binding:""`
	Title       string                `json:"title" binding:""`
	Description string                `json:"description" binding:""`
}

type VehicleModelUpdateInput struct {
	Category    *model.VehicleCategory `json:"category" binding:"required"`
	Label       *string                `json:"label" binding:""`
	IconImage   *string                `json:"iconImage" binding:""`
	CoverImage  *string                `json:"coverImage" binding:""`
	Title       *string                `json:"title" binding:""`
	Description *string                `json:"description" binding:""`
}
