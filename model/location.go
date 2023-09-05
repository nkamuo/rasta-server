package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Location struct {
	ID          uuid.UUID            `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Name        *string              `gorm:"type:varchar(255);" json:"name,omitempty"`
	Street      string               `gorm:"type:varchar(64)" json:"street,omitempty"`
	City        string               `gorm:"type:varchar(64)" json:"city,omitempty"`
	State       string               `gorm:"type:varchar(64)" json:"state,omitempty"`
	Country     string               `gorm:"type:varchar(64)" json:"country,omitempty"`
	CityCode    string               `gorm:"type:varchar(64)" json:"cityCode,omitempty"`
	StateCode   string               `gorm:"type:varchar(64)" json:"stateCode,omitempty"`
	CountryCode string               `gorm:"type:varchar(64)" json:"countryCode,omitempty"`
	Address     string               `gorm:"type:varchar(255)" json:"address,omitempty"`
	Coordinates *LocationCoordinates `gorm:"embedded;columnPrefix:'coords_'" json:"coords"`
	PostCode    string               `json:"postcode"`
	GoogleID    *string              `gorm:"type:varchar(64);" json:"googleId,omitempty"`
	Description *string              `gorm:"" json:"description,omitempty"`
	CreatedAt   time.Time            `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time            `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

type LocationCoordinates struct {
	Latitude  float32  `json:"latitude,omitempty"`
	Longitude float32  `json:"longitude,omitempty"`
	Altitude  float32  `json:"altitude,omitempty"`
	Accuracy  *float32 `json:"accuracy,omitempty"`
	Heading   *float32 `json:"heading,omitempty"`
	Speed     *float32 `json:"speed,omitempty"`
}

func (location *Location) BeforeCreate(tx *gorm.DB) (err error) {
	location.ID = uuid.New()
	location.CreatedAt = time.Now()
	location.UpdatedAt = time.Now()
	return nil
}

func (location Location) GetReference() (ref string) {
	if location.GoogleID != nil {
		return fmt.Sprintf("place_id:%s", *location.GoogleID)
	}
	if location.Coordinates != nil {
		coords := location.Coordinates
		return fmt.Sprintf("%f,%f", coords.Latitude, coords.Longitude)
	}
	return location.Address
}

func CreateLocationFromCoordinates(latitude, longitude float32) (location *Location, err error) {
	return &Location{
		Coordinates: &LocationCoordinates{
			Latitude:  latitude,
			Longitude: longitude,
		},
	}, nil

}
