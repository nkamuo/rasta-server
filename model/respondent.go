package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Respondent struct {
	// -
	ID        uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Active    bool      `gorm:"default:false;not null" json:"active"`
	Published bool      `gorm:"default:false;not null" json:"published"`

	//PRIMARY VEHICLE USED BY THIS RESPONDENT
	VehicleID *uuid.UUID `gorm:"unique" json:"vehicleId,omitempty"`
	Vehicle   *Vehicle   `gorm:"foreignKey:VehicleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"vehicle,omitempty"`

	//ASSOCIATED USER ACCOUNT
	UserID *uuid.UUID `gorm:";unique" json:"userId,omitempty"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	//COMPANY OF THE REPSONDANT
	CompanyID *uuid.UUID `gorm:"" json:"companyId,omitempty"`
	Company   *Company   `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"company,omitempty"`
	// PLACE THE RESPONDENT CAN OPERATE IN
	PlaceID *uuid.UUID `gorm:"unique" json:"placeId,omitempty"`
	Place   *Place     `gorm:"foreignKey:PlaceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"place,omitempty"`
	// TIMESTAMPS
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (Respondent *Respondent) BeforeCreate(tx *gorm.DB) (err error) {
	Respondent.ID = uuid.New()
	return nil
}
