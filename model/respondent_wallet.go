package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RespondentWallet struct {
	Wallet
	//THE RESPONDENT Who own the wallet
	RespondentID *uuid.UUID  `gorm:"not null" json:"respondentId,omitempty"`
	Respondent   *Respondent `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"`
}

func (wallet *RespondentWallet) BeforeCreate(tx *gorm.DB) (err error) {
	wallet.ID = uuid.New()
	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()
	return nil
}
