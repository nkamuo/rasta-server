package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItemVehicleInfo struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`

	ModelID *uuid.UUID    `gorm:"" json:"modelId,omitempty"`
	Model   *VehicleModel `gorm:"foreignKey:ModelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"model,omitempty"`

	Description *string   `gorm:"" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (orderItemVehicleInfo *OrderItemVehicleInfo) BeforeCreate(tx *gorm.DB) (err error) {
	orderItemVehicleInfo.ID = uuid.New()
	return nil
}

type OrderItemFuelTypeInfo struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`

	FuelTypeCode string     `gorm:"not null" json:"fuelTypeCode,omitempty"`
	FuelTypeID   *uuid.UUID `gorm:"" json:"fuelTypeId,omitempty"`
	FuelType     *Place     `gorm:"foreignKey:FuelTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"fuelType,omitempty"`
	//
	Description *string   `gorm:"" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (orderItemFuelTypeInfo *OrderItemFuelTypeInfo) BeforeCreate(tx *gorm.DB) (err error) {
	orderItemFuelTypeInfo.ID = uuid.New()
	return nil
}
