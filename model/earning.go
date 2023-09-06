package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderEarning struct {
	ID          uuid.UUID          `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Amount      uint64             `gorm:"" json:"amount"`
	Label       string             `gorm:"type:varchar(128);not null" json:"label,omitempty"`
	Description string             `gorm:"not null" json:"description,omitempty"`
	Status      OrderEarningStatus `gorm:"type:varchar(16);not null" json:"status,omitempty"`
	//THE RESPONDENT EXECUTING THIS TASK
	RequestID *uuid.UUID `gorm:"not null" json:"requestId,omitempty"`
	Request   *Request   `gorm:"foreignKey:RequestID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"request,omitempty"`

	//TIMESTAMPs
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

type OrderEarningStatus = string

const ORDER_EARNING_STATUS_PENDING = "pending"
const ORDER_EARNING_STATUS_COMPLETED = "completed"
const ORDER_EARNING_STATUS_CANCELLED = "cancelled"

func (earning OrderEarning) IsCommited() bool {
	return earning.Status == ORDER_EARNING_STATUS_COMPLETED
}

func (earning OrderEarning) IsPending() bool {
	return earning.Status == ORDER_EARNING_STATUS_PENDING
}
