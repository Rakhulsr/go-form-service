package domain

import (
	"time"

	"gorm.io/gorm"
)

type Form struct {
	gorm.Model
	Title       string     `gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:text"`
	UserID      uint       `gorm:"not null"`
	Question    []Question `gorm:"constraint:OnDelete:CASCADE;"`
}

type Question struct {
	ID       uint     `gorm:"primaryKey"`
	FormID   uint     `gorm:"not null"`
	Type     string   `gorm:"type:varchar(50);not null"`
	Text     string   `gorm:"type:text;not null"`
	Required bool     `gorm:"default:false"`
	Options  []Option `gorm:"constraint:OnDelete:CASCADE;"`
}

type Option struct {
	ID         uint   `gorm:"primaryKey"`
	QuestionID uint   `gorm:"not null"`
	Text       string `gorm:"type:varchar(255);not null"`
}

type Response struct {
	ID        uint `gorm:"primaryKey"`
	FormID    uint `gorm:"not null"`
	CreatedAt time.Time
	Answers   []Answer `gorm:"constraint:OnDelete:CASCADE;"`
}

type Answer struct {
	ID         uint `gorm:"primaryKey"`
	ResponseID uint `gorm:"not null"`
	QuestionID uint `gorm:"not null"`
	OptionID   *uint
	Text       *string
}
