package migrations

import (
	"github.com/Rakhulsr/go-form-service/internal/models/domain"
	"gorm.io/gorm"
)

func AutoMigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(&domain.User{}, &domain.Session{})
}
