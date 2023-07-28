package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Respondent struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Active    bool      `gorm:"default:false;not null" json:"active"`
	Published bool      `gorm:"default:false;not null" json:"published"`
	//ASSOCIATED USER ACCOUNT
	UserID *uuid.UUID `gorm:"unique" json:"userId,omitempty"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"` // BelongsToMany Association - These mappings may be optional
	//COMPANY OF THE REPSONDANT
	CompanyID uuid.UUID `gorm:"unique" json:"companyId,omitempty"`
	Company   *Company  `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"company,omitempty"` // BelongsToMany Association - These mappings may be optional
	// PLACE THIS RESPONDENT IS LOCALIZED
	PlaceID uuid.UUID `gorm:"unique" json:"placeId,omitempty"`
	Place   *Place    `gorm:"foreignKey:PlaceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"place,omitempty"` // BelongsToMany Association - These mappings may be optional
	// TIMESTAMPS
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (Respondent *Respondent) BeforeCreate(tx *gorm.DB) (err error) {
	Respondent.ID = uuid.New()
	return nil
}

// BeforeCreate will set a UUID rather than numeric ID.
// func (base *Base) BeforeCreate(scope *gorm.Scope) error {
// 	uuid, err := uuid.NewV4()
// 	if err != nil {
// 	 return err
// 	}
// 	return scope.SetColumn("ID", uuid)
//    }
