package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderAdjustment struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Amount      int64     `json:"amount,omitempty"`
	Code        string    `gorm:"not null;varchar(32)" json:"code,omitempty"`
	Title       string    `gorm:"not null" json:"title,omitempty"`
	Description *string   `gorm:"" json:"description,omitempty"`

	//THE MAIN ORDER ENTITY
	OrderID *uuid.UUID `gorm:"not null" json:"orderId,omitempty"`
	Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order,omitempty"`
}

func (adjustment *OrderAdjustment) BeforeCreate(tx *gorm.DB) (err error) {
	adjustment.ID = uuid.New()
	return nil
}
