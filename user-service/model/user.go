package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:32;not null"`
	Password string `gorm:"size:128;not null"`
	Email    string `gorm:"size:128;not null"`
}
