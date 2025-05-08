package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `gorm:"size:128;not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"not null"`
	MainImage   string  `gorm:"size:256"`
}
