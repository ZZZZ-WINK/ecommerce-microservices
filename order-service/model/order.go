package model

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     int64       `gorm:"not null;index"`
	TotalPrice float64     `gorm:"not null"`
	Status     int         `gorm:"not null"` // 对应proto中的OrderStatus枚举
	Items      []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID     uint    `gorm:"not null;index"`
	ProductID   int64   `gorm:"not null"`
	ProductName string  `gorm:"size:128;not null"`
	Quantity    int     `gorm:"not null"`
	Price       float64 `gorm:"not null"`
}
