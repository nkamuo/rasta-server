package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItem struct {
	ID               uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Rate             uint64    `gorm:"" json:"rate"`
	Quantity         uint64    `json:"quantity,omitempty"`
	AdjustmentsTotal uint64    `json:"adjustmentTotal,omitempty"`
	Total            uint64    `json:"total,omitempty"`

	//THE MAIN ORDER ENTITY
	OrderID *uuid.UUID `gorm:"not null" json:"orderId,omitempty"`
	Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order,omitempty"`

	OriginID *uuid.UUID `json:"originId"`
	Origin   *Location  `gorm:"" json:"origin,omitempty"`

	DestinationID *uuid.UUID `gorm:"not null" json:"destinationId"`
	Destination   *Location  `gorm:"" json:"destination,omitempty"`

	//ASSOCIATED PRODUCT
	ProductID *uuid.UUID `gorm:"" json:"productId,omitempty"`
	Product   *Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product,omitempty"`

	// OPTIONAL VEHICLE INFORMATIOn PROVIDED FOR ORDERS WHERE VEHICLE INFORMATION IS NECCESARY
	VehicleInfoID *uuid.UUID            `gorm:"" json:"vehicleInfoId,omitempty"`
	VehicleInfo   *OrderItemVehicleInfo `gorm:"foreignKey:FuelTypeInfoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"vehicleInfo,omitempty"`
	// OPTIONAL FUEL TYPE INFORMATION TO DESCRIBE FUEL TYPE NEEDED BY THE REQUESTING USER
	FuelTypeInfoID *uuid.UUID             `gorm:"" json:"fuelTypeInfoId,omitempty"`
	FuelTypeInfo   *OrderItemFuelTypeInfo `gorm:"foreignKey:FuelTypeInfoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"fuelType,omitempty"`

	//TIMESTAMPs
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (orderItem *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	orderItem.ID = uuid.New()
	return nil
}
