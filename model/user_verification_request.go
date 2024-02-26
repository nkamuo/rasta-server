package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserVerificationRequest struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Email       string    `gorm:"type:varchar(128);not null" json:"email,omitempty"`
	Code        string    `gorm:"type:varchar(255);not null" json:"code,omitempty"`
	Token       string    `gorm:"type:varchar(255);not null" json:"token,omitempty"`
	RequestType string    `gorm:"type:varchar(255);not null" json:"requestType,omitempty"`
	ExpiresAt   time.Time `gorm:"not null;" json:"expiresAt,omitempty"`
	CreatedAt   time.Time `gorm:"not null;" json:"createdAt,omitempty"`
	UpdatedAt   time.Time `gorm:"not null;ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (request *UserVerificationRequest) BeforeCreate(tx *gorm.DB) (err error) {
	request.ID = uuid.New()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	return nil
}
