package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
	Password  string `json:"-"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}
