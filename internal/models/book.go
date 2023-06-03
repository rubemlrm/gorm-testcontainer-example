package models

import (
	"time"
)

type Book struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
