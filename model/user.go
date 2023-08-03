package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Email          string    `gorm:"type:varchar(64);uniqueIndex:idx_notes_email,LENGTH(255);not null" json:"email,omitempty"`
	Phone          string    `gorm:"type:varchar(64);uniqueIndex:idx_notes_phone,LENGTH(255);not null" json:"phone,omitempty"`
	HashedPassword string    `gorm:"not null" json:"hashedPassword,omit"`
	FirstName      string    `gorm:"varchar(100)" json:"firstName,omitempty"`
	LastName       string    `gorm:"varchar(100)" json:"lastName,omitempty"`
	IsAdmin        *bool     `gorm:"default:false;not null" json:"isAdmin"`
	Published      bool      `gorm:"default:false;not null" json:"published"`
	CreatedAt      time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt      time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (User *User) BeforeCreate(tx *gorm.DB) (err error) {
	User.ID = uuid.New()
	return nil
}

func (user User) FullName() (fullname string) {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func (u *User) PrepareGive() {
	u.HashedPassword = ""
}

// BeforeCreate will set a UUID rather than numeric ID.
// func (base *Base) BeforeCreate(scope *gorm.Scope) error {
// 	uuid, err := uuid.NewV4()
// 	if err != nil {
// 	 return err
// 	}
// 	return scope.SetColumn("ID", uuid)
//    }
