	package services

	import (
		"errors"
		"fmt"
		"log"
		"time"

		"github.com/Rakhulsr/go-form-service/internal/models/domain"
		"github.com/Rakhulsr/go-form-service/internal/repositories"
		"github.com/Rakhulsr/go-form-service/utils"
		"github.com/go-playground/validator/v10"
		"gorm.io/gorm"
	)

	type UserService interface {
		FindByID(userID int) (*domain.User, error)
		FindOrCreateByEmail(email string, googleID string) (*domain.User, error)
		LoginGoogle(email, googleID string, userID uint) (string, string, error)
		RefreshSession(userID uint) (*domain.Session, error)
	}

	type UserServiceImpl struct {
		UserRepo repositories.UserRepositories
		Validate *validator.Validate
	}

	func NewUserServiceImpl(userRepo repositories.UserRepositories, validate *validator.Validate) *UserServiceImpl {
		return &UserServiceImpl{UserRepo: userRepo, Validate: validate}
	}

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

	func (s *UserServiceImpl) LoginGoogle(email, googleID string, userID uint) (string, string, error) {
		user, err := s.FindOrCreateByEmail(email, googleID)
		if err != nil {
			return "", "", err
		}

		accessToken, err := utils.GenerateAccesToken(user.Email, user.GoogleID, userID)
		if err != nil {
			return "", "", err
		}

		refreshToken, err := utils.GenerateRefreshToken(email, googleID, userID)
		if err != nil {
			return "", "", err
		}

		session := domain.Session{
			UserID:    user.ID,
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}

		if err := s.UserRepo.CreateSession(session); err != nil {
			return "", "", err
		}

		return accessToken, refreshToken, nil
	}

	func (s *UserServiceImpl) RefreshSession(userID uint) (*domain.Session, error) {

		session, err := s.UserRepo.FindSessionByUserID(userID)
		if err != nil {
			log.Println("Session expired or invalid:", err)
			return nil, errors.New("unauthorized: session expired")
		}

		return session, nil

	}
