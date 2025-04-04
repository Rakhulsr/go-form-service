package repositories

import (
	"github.com/Rakhulsr/go-form-service/internal/models/domain"
	"gorm.io/gorm"
)

type FormRepositories struct {
	Db *gorm.DB
}

func NewFormRepository(db *gorm.DB) *FormRepositories {
	return &FormRepositories{Db: db}
}

func (r *FormRepositories) CreateForm(form domain.Form) error {
	return r.Db.Create(form).Error
}

func (r *FormRepositories) GetForms() (*[]domain.Form, error) {
	var forms []domain.Form
	err := r.Db.Preload("Questions.Options").Find(&forms).Error
	return &forms, err
}

func (r *FormRepositories) GetFormByID(id uint) (*domain.Form, error) {
	var form domain.Form
	err := r.Db.Preload("Questions.Options").First(&form, id).Error
	return &form, err
}

func (r *FormRepositories) GetFormsByUserID(userID uint) (*[]domain.Form, error) {
	var forms []domain.Form

	if err := r.Db.Where("user_id=?", userID).Find(&forms).Error; err != nil {
		return nil, err
	}

	return &forms, nil

}
