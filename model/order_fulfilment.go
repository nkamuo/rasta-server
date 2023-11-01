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

	//The Session of the responder when accepting the request
	SessionID *uuid.UUID         `gorm:"" json:"sessionId,omitempty"`
	Session   *RespondentSession `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"session,omitempty"`

	//This will be set manually from the gin-controller to let the end client tell that the responder in in motion
	Coordinates *LocationCoordinates `gorm:"embedded;columnPrefix:'coords_'" json:"coords,omitempty"`

	// IDENTIY AND VEHICLE CONFIRMATION
	VerifiedResponderAt *time.Time `gorm:"" json:"verifiedResponderAt,"`
	VerifiedClientAt    *time.Time `gorm:"" json:"verifiedClientAt,"`
	//
	ResponderConfirmedAt *time.Time `gorm:"" json:"ResponderConfirmedAt,"`
	ClientConfirmedAt    *time.Time `gorm:"" json:"clientConfirmedAt,"`
	AutoConfirmedAt      *time.Time `gorm:"" json:"autoConfirmedAt,"` //THIS IS AUTOMATICALLY SET IF THE CLIENT DOES NOT CONFIRM ON TIME
	//

	InitialExpectedTimeOfAt *time.Time `gorm:"" json:"initialExpectedTimeOfAt,"` //THIS IS SET ONCE THE CLIENT ACCEPTS THE REQUST
	ExpectedTimeOfAt        *time.Time `gorm:"" json:"expectedTimeOfAt,"`        // COMPUTED PER RESPONDENT REQUEST UPDATE REQUEST
	//
	CreatedAt time.Time `gorm:"not null;" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (fulfilment *OrderFulfilment) BeforeCreate(tx *gorm.DB) (err error) {
	fulfilment.ID = uuid.New()
	fulfilment.CreatedAt = time.Now()
	fulfilment.UpdatedAt = time.Now()
	return nil
}

func (fulfilment *OrderFulfilment) IsComplete() (isComplete bool) {
	if fulfilment.ClientConfirmedAt != nil || fulfilment.AutoConfirmedAt != nil {
		return true
	}
	return false
}
