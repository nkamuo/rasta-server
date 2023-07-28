package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID       `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Title       string          `gorm:"type:varchar(255);uniqueIndex:idx_notes_title,LENGTH(255);not null" json:"title,omitempty"`
	Description string          `gorm:"not null" json:"description,omitempty"`
	Category    ProductCategory `gorm:"varchar(100)" json:"category,omitempty"`
	PlaceID     uuid.UUID       `gorm:"unique" json:"companyId,omitempty"`
	Place       *Place          `gorm:"foreignKey:PlaceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"company"` // BelongsToMany Association - These mappings may be optional
	Published   bool            `gorm:"default:false;not null" json:"published"`
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

// BeforeCreate will set a UUID rather than numeric ID.
// func (base *Base) BeforeCreate(scope *gorm.Scope) error {
// 	uuid, err := uuid.NewV4()
// 	if err != nil {
// 	 return err
// 	}
// 	return scope.SetColumn("ID", uuid)
//    }
