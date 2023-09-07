package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RespondentEarning struct {
	OrderEarning
	//THE RESPONDENT EXECUTING THIS TASK
	RespondentID *uuid.UUID  `gorm:"not null" json:"respondentId,omitempty"`
	Respondent   *Respondent `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"`
}

func (earning *RespondentEarning) BeforeCreate(tx *gorm.DB) (err error) {
	earning.ID = uuid.New()
	earning.CreatedAt = time.Now()
	earning.UpdatedAt = time.Now()
	return nil
}
