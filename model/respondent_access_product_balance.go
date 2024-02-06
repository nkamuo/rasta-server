package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RespondentAccessProductBalance struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`

	//ASSOCIATED USER ACCOUNT
	RespondentID *uuid.UUID  `gorm:"unique;not null" json:"respondentId,omitempty"`
	Respondent   *Respondent `gorm:"foreignKey:RespondentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"respondent,omitempty"` // BelongsToMany Association - These mappings may be optional
	//Balance
	Balance *int64 `gorm:"default:0;not null" json:"balance"`
	// The total number of request the user has paid for and can handle
	//TODO: Decremete this number each time he handles a request

	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (balance *RespondentAccessProductBalance) Increment(quantity int64) {
	if balance.Balance == nil {
		var Val int64 = 0
		balance.Balance = &Val
	}
	*balance.Balance += quantity
}

func (balance *RespondentAccessProductBalance) BeforeCreate(tx *gorm.DB) (err error) {
	balance.ID = uuid.New()
	balance.CreatedAt = time.Now()
	balance.UpdatedAt = time.Now()

	return nil
}
