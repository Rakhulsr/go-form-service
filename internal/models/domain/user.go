package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	GoogleID string `gorm:"type:varchar(255);uniqueIndex" json:"google_id"`
}
