package dto

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	Name      string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"omitempty"`
}