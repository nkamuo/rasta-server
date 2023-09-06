package model

import (
	"time"

	"github.com/google/uuid"
)

type Withdrawal struct {
	ID          uuid.UUID        `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Amount      uint64           `gorm:"" json:"amount"`
	Description string           `gorm:"type:varchar(255)" json:"description,omitempty"`
	Status      WithdrawalStatus `gorm:"type:varchar(16);not null" json:"status,omitempty"`
	// //THE RESPONDENT EXECUTING THIS TASK
	// WithdrawalMethodID *uuid.UUID `gorm:"not null" json:"requestId,omitempty"`
	// WithdrawalMethod   *WithdrawalMethod   `gorm:"foreignKey:WithdrawalMethodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"request,omitempty"`

	StripePayoutID string `gorm:"type:varchar(255)" json:"stripePayoutID,omitempty"`
	//WHEN THE WITHDRAWAL IS COMMITED TO THE CLIENTS ACCOUNT
	ComittedAt *time.Time `gorm:"" json:"comittedAt,omitempty"`

	//TIMESTAMPs
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (withdrawal Withdrawal) IsCommited() bool {
	return withdrawal.Status == ORDER_WITHDRAWAL_STATUS_COMPLETED
}

func (withdrawal Withdrawal) IsPending() bool {
	return withdrawal.Status == ORDER_WITHDRAWAL_STATUS_PENDING
}

type WithdrawalStatus = string

const ORDER_WITHDRAWAL_STATUS_PENDING = "pending"
const ORDER_WITHDRAWAL_STATUS_COMPLETED = "completed"
const ORDER_WITHDRAWAL_STATUS_CANCELLED = "cancelled"
