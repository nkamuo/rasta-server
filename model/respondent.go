package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Respondent struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	// Title       string    `gorm:"type:varchar(255);uniqueIndex:idx_notes_title,LENGTH(255);not null" json:"title,omitempty"`
	// Description string    `gorm:"not null" json:"description,omitempty"`
	// Category    string    `gorm:"varchar(100)" json:"category,omitempty"`
	Active    bool       `gorm:"default:false;not null" json:"active"`
	Published bool       `gorm:"default:false;not null" json:"published"`
	UserID    *uuid.UUID `gorm:"unique" json:"userId,omitempty"`
	User      *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"` // BelongsToMany Association - These mappings may be optional
	//COMPANY OF THE REPSONDANT
	CompanyID uuid.UUID `gorm:"unique" json:"companyId,omitempty"`
	Company   *Company  `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"company"` // BelongsToMany Association - These mappings may be optional
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
