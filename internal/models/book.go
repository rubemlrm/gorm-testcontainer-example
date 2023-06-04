package models

import (
	"time"
)

type Book struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"not null;type:varchar(12)"`
	Author    string    `gorm:"not null;type:varchar(12)"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
}
