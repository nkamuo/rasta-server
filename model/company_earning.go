package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyEarning struct {
	OrderEarning
	//THE RESPONDENT EXECUTING THIS TASK
	CompanyID *uuid.UUID `gorm:"not null" json:"companyId,omitempty"`
	Company   *Company   `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"company,omitempty"`
	//THE RESPONDENT THAT DID THE JOB
	//THE RESPONDENT EXECUTING THIS TASK
	RespondentID *uuid.UUID  `gorm:"not null" json:"respondentId,omitempty"`
	Respondent   *Respondent `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"`
}

func (earning *CompanyEarning) BeforeCreate(tx *gorm.DB) (err error) {
	earning.ID = uuid.New()
	earning.CreatedAt = time.Now()
	earning.UpdatedAt = time.Now()
	return nil
}
