package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TowingPlaceRate struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	//
	PlaceID *uuid.UUID `gorm:"not null;" json:"placeId,omitempty"`
	Place   *Place     `gorm:"foreignKey:PlaceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"place,omitempty"`
	//
	// Sequence *uint32 `gorm:"not null;uniqueIndex:idx_place_sequence" json:"sequence"`
	MinDistance *uint64 `gorm:"not null;" json:"minDistance"`
	MaxDistance *uint64 `gorm:"not null;" json:"maxDistance"`
	//
	Rate *uint64 `gorm:"not null;" json:"rate"`
	//
	Description *string   `gorm:"" json:"description,omitempty"`
	Active      *bool     `gorm:"default:false;not null" json:"active"`
	CreatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (towingPlaceRate *TowingPlaceRate) BeforeCreate(tx *gorm.DB) (err error) {
	towingPlaceRate.ID = uuid.New()
	return nil
}
