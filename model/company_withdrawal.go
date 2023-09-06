package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyWithdrawal struct {
	Withdrawal
	//THE RESPONDENT EXECUTING THIS TASK
	// CompanyID *uuid.UUID `gorm:"not null" json:"companyId,omitempty"`
	// Company   *Company   `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"company,omitempty"`
	//THE RESPONDENT THAT DID THE JOB
	//THE RESPONDENT EXECUTING THIS TASK
	WalletID *uuid.UUID     `gorm:"not null" json:"walletId,omitempty"`
	Wallet   *CompanyWallet `gorm:"foreignKey:WalletID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"wallet,omitempty"`
}

func (earning *CompanyWithdrawal) BeforeCreate(tx *gorm.DB) (err error) {
	earning.ID = uuid.New()
	earning.CreatedAt = time.Now()
	earning.UpdatedAt = time.Now()
	return nil
}
