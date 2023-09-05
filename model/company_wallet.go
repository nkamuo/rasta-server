package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyWallet struct {
	Wallet
	//THE Company Owner of this Wallet
	CompanyID *uuid.UUID `gorm:"not null;unique" json:"companyId,omitempty"`
	Company   *Company   `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"company,omitempty"`
}

func (wallet *CompanyWallet) BeforeCreate(tx *gorm.DB) (err error) {
	wallet.ID = uuid.New()
	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()
	return nil
}
