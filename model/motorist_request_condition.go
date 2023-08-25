package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MotoristRequestSituation struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	//
	Title     string `gorm:"" json:"title,omitempty"`
	SubTitlte string `gorm:"" json:"subtitle,omitempty"`
	//
	Code string `gorm:"" json:"code,omitempty"`
	//
	Description string `gorm:"" json:"description,omitempty"`
	Note        string `gorm:"" json:"note,omitempty"`
	//
	// Active      *bool     `gorm:"default:false;not null" json:"active"`
	// CreatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	// UpdatedAt   time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (situation *MotoristRequestSituation) BeforeCreate(tx *gorm.DB) (err error) {
	situation.ID = uuid.New()
	return nil
}

var defaultConfitions = []MotoristRequestSituation{
	{
		Code:        "in_a_dark_place",
		Title:       "In a Dark place",
		SubTitlte:   "In a Dark Place",
		Description: "You are in a dark place",
		Note:        "The motorist is in a dark place",
	},
	{
		Code:        "on_high_way",
		Title:       "On the highway",
		SubTitlte:   "Currently On the high way",
		Description: "You are on the high way",
		Note:        "The motorist is on the highway",
	},
}

func GetDefaultMotoristRequestSituations() (situations []MotoristRequestSituation) {
	return defaultConfitions
}
