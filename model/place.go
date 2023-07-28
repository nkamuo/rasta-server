package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Place struct {
	ID          uuid.UUID     `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Name        string        `gorm:"type:varchar(255);uniqueIndex:idx_place_name,LENGTH(255);not null" json:"name,omitempty"`
	ShortName   string        `gorm:"type:varchar(255);uniqueIndex:idx_place_short_name,LENGTH(255);" json:"shortName,omitempty"`
	LongName    string        `gorm:"type:varchar(255);uniqueIndex:idx_place_long_name,LENGTH(255);" json:"longName,omitempty"`
	Description string        `gorm:"" json:"description,omitempty"`
	Category    PlaceCategory `gorm:"varchar(100)" json:"category,omitempty"`
	Active      bool          `gorm:"default:false;not null" json:"active"`
	CreatedAt   time.Time     `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time     `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (company *Place) BeforeCreate(tx *gorm.DB) (err error) {
	company.ID = uuid.New()
	return nil
}

type PlaceCategory = string

const (
	PLACE_CITY    PlaceCategory = "CITY"
	PLACE_STATE   PlaceCategory = "STATE"
	PLACE_COUNTRY PlaceCategory = "COUNTRY"
)
