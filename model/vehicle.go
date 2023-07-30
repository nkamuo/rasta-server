package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	ID                 uuid.UUID     `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Description        string        `gorm:"not null" json:"description,omitempty"`
	LicensePlaceNumber string        `gorm:"" json:"licensePlaceNumber,omitempty"`
	OwnerID            uuid.UUID     `gorm:"not null" json:"ownerId,omitempty"`
	Owner              *User         `gorm:"" json:"owner"`
	ModelID            uuid.UUID     `gorm:"not null" json:"modelId,omitempty"`
	Model              *VehicleModel `gorm:"" json:"model"`
	Published          bool          `gorm:"default:false;not null" json:"published"`
	CreatedAt          time.Time     `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt          time.Time     `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (product *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return nil
}
