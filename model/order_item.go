package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID               uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	UnitPrice        int64     `gorm:"" json:"unitPrice"`
	Quantity         int64     `json:"quantity,omitempty"`
	AdjustmentsTotal int64     `json:"adjustmentTotal,omitempty"`
	Total            int64     `json:"total,omitempty"`

	//THE MAIN ORDER ENTITY
	OrderID *uuid.UUID `gorm:"not null" json:"orderId,omitempty"`
	Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"order,omitempty"`

	OriginID *uuid.UUID `json:"originId"`
	Origin   *Location  `gorm:"" json:"origin,omitempty"`

	DestinationID *uuid.UUID `gorm:"not null" json:"destinationId"`
	Destination   *Location  `gorm:"" json:"destination,omitempty"`

	//ASSOCIATED USER ACCOUNT
	ProductID *uuid.UUID `gorm:"not null" json:"productId,omitempty"`
	Product   *Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product,omitempty"`

	//TIMESTAMPs
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}
