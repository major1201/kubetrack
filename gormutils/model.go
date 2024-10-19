package gormutils

import (
	"time"

	"gorm.io/gorm"
)

type ModelUnscoped struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Model base DB model
type Model struct {
	ModelUnscoped

	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
