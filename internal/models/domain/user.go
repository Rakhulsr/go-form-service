package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	GoogleID string `gorm:"type:varchar(255);uniqueIndex" json:"google_id"`
}

type Session struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Token     string    `gorm:"type:varchar(512);unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
