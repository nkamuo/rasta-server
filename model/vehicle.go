package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	ID                 uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	LicensePlateNumber string    `gorm:"" json:"licensePlateNumber,omitempty"`
	Color              string    `gorm:"varchar(64);" json:"color,omitempty"`
	// USER OWNER
	OwnerID *uuid.UUID `gorm:"" json:"ownerId,omitempty"`
	Owner   *User      `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"owner,omitempty"`

	// COMPANY OWNER
	CompanyID *uuid.UUID `gorm:"" json:"companyId,omitempty"`
	Company   *Company   `gorm:"" json:"company"`

	// THE Model Of the Vehicle -- Optional;
	ModelID *uuid.UUID    `gorm:"" json:"modelId,omitempty"`
	Model   *VehicleModel `gorm:"" json:"model"`
	//  ALTERNATIVE TO VEHICLE MODEL
	VinNumber *string `gorm:"" json:"vin,omitempty"`
	MakeName  *string `gorm:"" json:"makeName,omitempty"`
	ModelName *string `gorm:"" json:"modelName,omitempty"`
	//
	Published   *bool     `gorm:"default:false;not null" json:"published"`
	Description string    `gorm:"not null" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (product *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	return nil
}
