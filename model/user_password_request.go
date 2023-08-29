package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPasswordResetRequest struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	//
	UserID *uuid.UUID `gorm:"unique;not null;" json:"userId,omitempty"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	//
	Token string `gorm:"not null" json:"token,omit"`
	//
	ExpiresAt time.Time `gorm:"not null;" json:"expiresAt,omitempty"`
	CreatedAt time.Time `gorm:"not null;" json:"createdAt,omitempty"`
}

func (request *UserPasswordResetRequest) BeforeCreate(tx *gorm.DB) (err error) {
	request.ID = uuid.New()
	return nil
}
