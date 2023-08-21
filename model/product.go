package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID       `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Label       string          `gorm:"type:varchar(128);not null" json:"label,omitempty"`
	Title       string          `gorm:"type:varchar(255);not null" json:"title,omitempty"`
	Description string          `gorm:"not null" json:"description,omitempty"`
	Category    ProductCategory `gorm:"type:varchar(100);not null;uniqueIndex:UNIQUE_PLACE_PRODUCT_CATEGORY" json:"category,omitempty"`
	PlaceID     uuid.UUID       `gorm:"uniqueIndex:UNIQUE_PLACE_PRODUCT_CATEGORY" json:"placeId,omitempty"`
	Place       *Place          `gorm:"foreignKey:PlaceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"place,omitempty"` // BelongsToMany Association - These mappings may be optional
	Rate        uint64          `gorm:"" json:"rate,omitempty"`
	Bundled     *bool           `gorm:"default:false;not null" json:"bundled"` //Tells that the rate specified when ordering this item applies for all the quantity and not each
	IconImage   string          `gorm:"" json:"iconImage,omitempty"`
	CoverImage  string          `gorm:"" json:"coverImage,omitempty"`
	Published   *bool           `gorm:"default:false;not null" json:"published"`
	CreatedAt   time.Time       `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time       `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return nil
}

type ProductCategory = string

const (
	PRODUCT_TOWING_SERVICE        ProductCategory = "TOWING_SERVICE"
	PRODUCT_FLAT_TIRE_SERVICE     ProductCategory = "FLAT_TIRE_SERVICE"
	PRODUCT_FUEL_DELIVERY_SERVICE ProductCategory = "FUEL_DELIVERY_SERVICE"
	PRODUCT_JUMP_START_SERVICE    ProductCategory = "JUMP_START_SERVICE"
	PRODUCT_KEY_UNLOCK_SERVICE    ProductCategory = "KEY_UNLOCK_SERVICE"
	PRODUCT_TIRE_AIR_SERVICE      ProductCategory = "TIRE_AIR_SERVICE"
)

func ValidateProductCategory(category ProductCategory) (err error) {
	switch category {
	case PRODUCT_FLAT_TIRE_SERVICE:
		return nil
	case PRODUCT_FUEL_DELIVERY_SERVICE:
		return nil
	case PRODUCT_TIRE_AIR_SERVICE:
		return nil
	case PRODUCT_TOWING_SERVICE:
		return nil
	case PRODUCT_JUMP_START_SERVICE:
		return nil
	case PRODUCT_KEY_UNLOCK_SERVICE:
		return nil
	}
	return errors.New(fmt.Sprintf("Unsupported Product Category \"%s\"", category))
}
