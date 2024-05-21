package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	PasswordHash string
	Email        string        `gorm:"unique"`
	Transactions []Transaction `gorm:"foreignKey:UserID"`
	Budgets      []Budget      `gorm:"foreignKey:UserID"`
}

type Transaction struct {
	gorm.Model
	UserID      uint
	Type        string
	Amount      float64
	Description string
	CategoryID  uint
	Category    Category
	Date        time.Time `gorm:"type:timestamp"`
}

type Category struct {
	gorm.Model
	Name     string `gorm:"unique"`
	IsActive bool   `gorm:"default:true"`
}

type Budget struct {
	gorm.Model
	UserID     uint
	Category   Category
	CategoryID uint
	Amount     float64
	StartDate  time.Time `gorm:"type:timestamp"`
	EndDate    time.Time `gorm:"type:timestamp"`
}
