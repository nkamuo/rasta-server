package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FuelType struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Code        string    `gorm:"type:varchar(6);not null" json:"code,omitempty"`
	Rate        uint64    `gorm:"" json:"rate"`
	Title       string    `gorm:"varchar(64)not null" json:"title,omitempty"`
	ShortName   *string   `gorm:"varchar(64)not null" json:"shortName,omitempty"`
	Description *string   `gorm:"" json:"description,omitempty"`
	Published   bool      `gorm:"default:false;not null" json:"published"`
	CreatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (fuelType *FuelType) BeforeCreate(tx *gorm.DB) (err error) {
	fuelType.ID = uuid.New()
	return nil
}

type FuelTypePlaceRate struct {
	ID   uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Rate uint64    `gorm:"" json:"rate"`
	//
	FuelTypeID *uuid.UUID `gorm:"not null" json:"fuelTypeId,omitempty"`
	FuelType   *FuelType  `gorm:"foreignKey:FuelTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"fuelType,omitempty"`
	//
	PlaceID *uuid.UUID `gorm:"not null" json:"placeId,omitempty"`
	Place   *Place     `gorm:"foreignKey:PlaceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"place,omitempty"`
	//
	Description *string   `gorm:"" json:"description,omitempty"`
	Active      *bool     `gorm:"default:false;not null" json:"active"`
	CreatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (fuelTypePlaceRate *FuelTypePlaceRate) BeforeCreate(tx *gorm.DB) (err error) {
	fuelTypePlaceRate.ID = uuid.New()
	return nil
}
