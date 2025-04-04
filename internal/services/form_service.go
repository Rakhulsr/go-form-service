package services

import (
	"errors"

	"github.com/Rakhulsr/go-form-service/internal/models/domain"
	"github.com/Rakhulsr/go-form-service/internal/repositories"
	"github.com/go-playground/validator/v10"
)

type FormService struct {
	FormRepo *repositories.FormRepositories
	Validate *validator.Validate
}

func NewFormService(repo *repositories.FormRepositories, validate *validator.Validate) *FormService {
	return &FormService{FormRepo: repo, Validate: validate}
}

func (s *FormService) CreateForm(form domain.Form) error {
	if form.Title == "" {
		return errors.New("Tittle is rqeuired")
	}

	return s.FormRepo.CreateForm(form)

}

func (s *FormService) GetForms() (*[]domain.Form, error) {
	return s.FormRepo.GetForms()
}

func (s *FormService) GetFormByID(id uint) (*domain.Form, error) {
	return s.FormRepo.GetFormByID(id)
}
