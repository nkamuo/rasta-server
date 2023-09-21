package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID                 uuid.UUID    `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Balance            uint64       `gorm:"not null" json:"balance"`
	ResultantBalance   *uint64      `json:"resultantBalance"`
	PendingCreditTotal uint64       `gorm:"not null" json:"pendingCreditTotal"`
	PendingDebitTotal  uint64       `gorm:"not null" json:"pendingDebitTotal"`
	Label              string       `gorm:"type:varchar(128);" json:"label,omitempty"`
	Description        string       `gorm:"not null" json:"description,omitempty"`
	Status             WalletStatus `gorm:"type:varchar(16);" json:"status,omitempty"`

	//TIMESTAMPs
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (wallet *Wallet) CalculateResultantBalance() (resultantBalance uint64) {
	resultantBalance = (wallet.Balance + wallet.PendingCreditTotal - wallet.PendingDebitTotal)
	wallet.ResultantBalance = &resultantBalance
	return resultantBalance
}

func (wallet *Wallet) InitEarning(earning OrderEarning) (err error) {
	if earning.Status != ORDER_EARNING_STATUS_PENDING {
		return errors.New("Cannot init non-pending earning")
	}
	wallet.PendingCreditTotal += earning.Amount
	return nil
}

func (wallet *Wallet) InitWithdrawal(withdrawal Withdrawal) (err error) {
	if withdrawal.Status != ORDER_WITHDRAWAL_STATUS_PENDING {
		return errors.New("Cannot init non-pending withdrawal")
	}
	if withdrawal.Amount > wallet.Balance {
		message := "Target amount is greater than account balance"
		return errors.New(message)
	}
	wallet.Balance += withdrawal.Amount
	// wallet.PendingDebitTotal += withdrawal.Amount
	return nil
}

func (wallet *Wallet) CommiteEarning(earning OrderEarning) (err error) {
	if earning.Status != ORDER_EARNING_STATUS_COMPLETED {
		return errors.New("Cannot commit incomplete earning")
	}
	wallet.PendingCreditTotal -= earning.Amount
	wallet.Balance += earning.Amount
	return nil
}

func (wallet *Wallet) CommiteWithdrawal(withdrawal Withdrawal) (err error) {
	if withdrawal.Status != ORDER_WITHDRAWAL_STATUS_COMPLETED {
		return errors.New("Cannot commit incomplete withdrawal")
	}
	wallet.PendingDebitTotal -= withdrawal.Amount
	wallet.Balance -= withdrawal.Amount
	return nil
}

type WalletStatus = string

const WALLET_STATUS_ACTIVE = "active"       // DEPOSIT AND WITHDRAWAL
const WALLET_STATUS_SUSPENDED = "suspended" // DEPOSIT BUT NO WITHDRAWAL
const WALLET_STATUS_CLOSED = "closed"       // NO DEPOSIT, NO WITHDRAWAL
