package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderFulfilment struct {
	ID     uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Status string    `gorm:"type:varchar(32);not null;" json:"status,omitempty"`
	// //THE MAIN ORDER ENTITY
	// OrderID *uuid.UUID `gorm:"not null" json:"orderId,omitempty"`
	// Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order,omitempty"`
	//
	//REQUESTING USER ACCOUNT
	ResponderID *uuid.UUID  `gorm:"" json:"responderId,omitempty"`
	Responder   *Respondent `gorm:"foreignKey:ResponderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"responder,omitempty"`

	// IDENTIY AND VEHICLE CONFIRMATION
	VerifiedResponderAt *time.Time `gorm:"" json:"verifiedResponderAt,"`
	VerifiedClientAt    *time.Time `gorm:"" json:"verifiedClientAt,"`
	//
	ResponderConfirmedAt *time.Time `gorm:"" json:"ResponderConfirmedAt,"`
	ClientConfirmedAt    *time.Time `gorm:"" json:"clientConfirmedAt,"`

	CreatedAt time.Time `gorm:"not null;" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (fulfilment *OrderFulfilment) BeforeCreate(tx *gorm.DB) (err error) {
	fulfilment.ID = uuid.New()
	fulfilment.CreatedAt = time.Now()
	fulfilment.UpdatedAt = time.Now()
	return nil
}
