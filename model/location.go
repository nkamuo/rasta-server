package model

import (
	"github.com/google/uuid"
)

type Location struct {
	ID          uuid.UUID            `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Name        string               `gorm:"type:varchar(255);uniqueIndex:idx_location_name,LENGTH(255);not null" json:"name,omitempty"`
	Street      string               `gorm:"type:varchar(255);uniqueIndex:idx_location_street,LENGTH(255);" json:"longName,omitempty"`
	Address     string               `gorm:"type:varchar(255);uniqueIndex:idx_location_address,LENGTH(255);" json:"shortName,omitempty"`
	Coordinates *LocationCoordinates `json:"coordinates"`
	GoogleID    *string              `gorm:"type:varchar(64);uniqueIndex:idx_location_google_id,LENGTH(64);not null" json:"googleId,omitempty"`
	Description *string              `gorm:"" json:"description,omitempty"`
	// CreatedAt   time.Time     `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	// UpdatedAt   time.Time     `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

type LocationCoordinates struct {
	Latitude  float32  `json:"latitude,omitempty"`
	Longitude float32  `json:"longitude,omitempty"`
	Altitude  float32  `json:"altitude,omitempty"`
	Accuracy  *float32 `json:"accuracy,omitempty"`
	Heading   *float32 `json:"heading,omitempty"`
	Speed     *float32 `json:"speed,omitempty"`
}
