package entity

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	Image       string    `json:"image" gorm:"type:varchar(255)"`
	CategoryId  uint      `json:"category_id"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
	// ID         uint      `json:"id" gorm:"primaryKey,autoIncrement"`
	// CreatedAt  time.Time `json:"created_at"`
	// UpdatedAt  time.Time `json:"updated_at"`
	// DeletedAt  gorm.DeletedAt `json:"-" gorm:"index,column:deleted_at"`
}