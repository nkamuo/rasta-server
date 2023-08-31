package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	ID             uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Title          string    `gorm:"type:varchar(255);uniqueIndex:idx_notes_title,LENGTH(255);not null" json:"title,omitempty"`
	LicenseNumber  string    `gorm:"varchar(64)" json:"licenseNumber,omitempty"`
	Description    string    `gorm:"" json:"description,omitempty"`
	Category       string    `gorm:"varchar(100)" json:"category,omitempty"`
	Active         *bool     `gorm:"default:false;not null" json:"active"`
	Published      *bool     `gorm:"default:false;not null" json:"published"`
	OperatorUserID uuid.UUID `gorm:"unique" json:"userId,omitempty"`
	OperatorUser   *User     `gorm:"foreignKey:OperatorUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"operatorUser"` // BelongsToMany Association - These mappings may be optional
	CreatedAt      time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt      time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (company *Company) BeforeCreate(tx *gorm.DB) (err error) {
	company.ID = uuid.New()
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()
	return nil
}
