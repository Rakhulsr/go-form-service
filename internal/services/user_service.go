package services

import (
	"errors"
	"fmt"

	"github.com/Rakhulsr/go-form-service/internal/models/domain"
	"github.com/Rakhulsr/go-form-service/internal/repositories"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserService interface {
	// FindByEmail(email string) (*domain.User, error)
	FindByID(userID int) (*domain.User, error)
	FindOrCreateByEmail(email string, googleID string) (*domain.User, error)
	// Register(user domain.User) error
}

type UserServiceImpl struct {
	UserRepo repositories.UserRepositories
	Validate *validator.Validate
}

func NewUserServiceImpl(userRepo repositories.UserRepositories, validate *validator.Validate) *UserServiceImpl {
	return &UserServiceImpl{UserRepo: userRepo, Validate: validate}
}

// func (s *UserServiceImpl) FindByEmail(email string) (*domain.User, error) {

// 	user, err := s.UserRepo.FindByEmail(email)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return user, nil
// }

func (s *UserServiceImpl) FindByID(userID int) (*domain.User, error) {
	user, err := s.UserRepo.FindById(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %v", err)
		}
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}
	return user, nil
}

func (s *UserServiceImpl) FindOrCreateByEmail(email string, googleID string) (*domain.User, error) {

	user, err := s.UserRepo.FindOrCreateByEmail(email, googleID)
	if err != nil {
		return nil, fmt.Errorf("error finding or creating user: %v", err)
	}
	return user, nil
}
