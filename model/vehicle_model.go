package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehicleModel struct {
	ID          uuid.UUID       `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Label       string          `gorm:"type:varchar(128);not null" json:"label,omitempty"`
	Title       string          `gorm:"type:varchar(255);not null" json:"title,omitempty"`
	Description string          `gorm:"not null" json:"description,omitempty"`
	Category    VehicleCategory `gorm:"" json:"category"`
	Published   bool            `gorm:"default:false;not null" json:"published"`
	IconImage   string          `gorm:"" json:"iconImage,omitempty"`
	CoverImage  string          `gorm:"" json:"coverImage,omitempty"`
	CreatedAt   time.Time       `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time       `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (product *VehicleModel) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return nil
}

type VehicleCategory = string

const (
	VEHICLE_PICKUP VehicleCategory = "PICKUP"
	VEHICLE_SEDAN  VehicleCategory = "SEDAN"
	VEHICLE_SUV    VehicleCategory = "SUV"
)

func ValidateVehicleCategory(category VehicleCategory) (err error) {
	switch category {
	case VEHICLE_PICKUP:
		return nil
	case VEHICLE_SEDAN:
		return nil
	case VEHICLE_SUV:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported Vehicle Category \"%s\"", category))
}
