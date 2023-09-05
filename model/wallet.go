package model

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID                 uuid.UUID    `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Balance            uint64       `gorm:"not null" json:"balance"`
	ResultantBalance   *uint64      `gorm:"not null" json:"resultantBalance"`
	PendingCreditTotal uint64       `gorm:"not null" json:"pendingCreditTotal"`
	PendingDebitTotal  uint64       `gorm:"not null" json:"pendingDebitTotal"`
	Label              string       `gorm:"type:varchar(128);not null" json:"label,omitempty"`
	Description        string       `gorm:"not null" json:"description,omitempty"`
	Status             WalletStatus `gorm:"type:varchar(16);not null" json:"status,omitempty"`

	//TIMESTAMPs
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (wallet *Wallet) CalculateResultantBalance() (resultantBalance uint64) {
	resultantBalance = (wallet.Balance + wallet.PendingCreditTotal - wallet.PendingDebitTotal)
	wallet.ResultantBalance = &resultantBalance
	return resultantBalance
}

type WalletStatus = string

const WALLET_STATUS_ACTIVE = "active"       // DEPOSIT AND WITHDRAWAL
const WALLET_STATUS_SUSPENDED = "suspended" // DEPOSIT BUT NO WITHDRAWAL
const WALLET_STATUS_CLOSED = "closed"       // NO DEPOSIT, NO WITHDRAWAL
