package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestVehicleInfo struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`

	ModelID *uuid.UUID    `gorm:"" json:"modelId,omitempty"`
	Model   *VehicleModel `gorm:"foreignKey:ModelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"model,omitempty"`

	MakeName           *string `gorm:"" json:"makeName,omitempty"`
	ModelName          *string `gorm:"" json:"modelName,omitempty"`
	BodyColor          *string `gorm:"" json:"bodyColor,omitempty"`
	LicensePlateNumber *string `gorm:"" json:"licensePlateNumber,omitempty"`

	Description *string   `gorm:"" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (RequestVehicleInfo *RequestVehicleInfo) BeforeCreate(tx *gorm.DB) (err error) {
	RequestVehicleInfo.ID = uuid.New()
	RequestVehicleInfo.CreatedAt = time.Now()
	RequestVehicleInfo.UpdatedAt = time.Now()
	return nil
}

type RequestFuelTypeInfo struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`

	FuelTypeCode string     `gorm:"not null" json:"fuelTypeCode,omitempty"`
	FuelTypeID   *uuid.UUID `gorm:"" json:"fuelTypeId,omitempty"`
	FuelType     *FuelType  `gorm:"foreignKey:FuelTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"fuelType,omitempty"`
	//
	Description *string   `gorm:"" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (RequestFuelTypeInfo *RequestFuelTypeInfo) BeforeCreate(tx *gorm.DB) (err error) {
	RequestFuelTypeInfo.ID = uuid.New()
	RequestFuelTypeInfo.CreatedAt = time.Now()
	RequestFuelTypeInfo.UpdatedAt = time.Now()
	return nil
}
