package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Place struct {
	ID          uuid.UUID        `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Code        string           `gorm:"type:varchar(64);uniqueIndex:idx_place_code,LENGTH(64);not null" json:"code,omitempty"`
	GoggleID    string           `gorm:"type:varchar(64);uniqueIndex:idx_place_goggle_id,LENGTH(64);not null" json:"googleId,omitempty"`
	Coordinates PlaceCoordinates `gorm:"embedded;columnPrefix:coords_" json:"coords"`
	Name        string           `gorm:"type:varchar(255);not null" json:"name,omitempty"`
	ShortName   string           `gorm:"type:varchar(255)" json:"shortName,omitempty"`
	LongName    string           `gorm:"type:varchar(255)" json:"longName,omitempty"`
	Description string           `gorm:"" json:"description,omitempty"`
	Category    PlaceCategory    `gorm:"varchar(100)" json:"category,omitempty"`
	Active      *bool            `gorm:"default:false;not null" json:"active"`
	CreatedAt   time.Time        `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time        `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (place *Place) BeforeCreate(tx *gorm.DB) (err error) {
	place.ID = uuid.New()
	place.CreatedAt = time.Now()
	place.UpdatedAt = time.Now()
	return nil
}

type PlaceCoordinates struct {
	Latitude  float32 `json:"latitude,omitempty" binding:"required"`
	Longitude float32 `json:"longitude,omitempty" binding:"required"`
}

type PlaceCategory = string

const (
	PLACE_CITY    PlaceCategory = "CITY"
	PLACE_STATE   PlaceCategory = "STATE"
	PLACE_COUNTRY PlaceCategory = "COUNTRY"
)
