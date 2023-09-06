package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RespondentWithdrawal struct {
	Withdrawal
	WalletID *uuid.UUID        `gorm:"not null" json:"walletId,omitempty"`
	Wallet   *RespondentWallet `gorm:"foreignKey:WalletID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"wallet,omitempty"`

	//THE RESPONDENT EXECUTING THIS TASK
	// RespondentID *uuid.UUID  `gorm:"not null;unique" json:"respondentId,omitempty"`
	// Respondent   *Respondent `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"`
}

func (withdrawal *RespondentWithdrawal) BeforeCreate(tx *gorm.DB) (err error) {
	withdrawal.ID = uuid.New()
	withdrawal.CreatedAt = time.Now()
	withdrawal.UpdatedAt = time.Now()
	return nil
}
