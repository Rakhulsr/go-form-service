package repositories

import (
	"errors"
	"log"

	"github.com/Rakhulsr/go-form-service/internal/models/domain"
	"gorm.io/gorm"
)

type UserRepositories interface {
	FindById(userId int) (*domain.User, error)
	CreateUser(user domain.User) error
	FindOrCreateByEmail(email string, googleID string) (*domain.User, error)
	CreateSession(token domain.Session) error
	FindSessionByUserID(userID uint) (*domain.Session, error)
}

type UserReporitoriesImpl struct {
	Db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) UserReporitoriesImpl {
	return UserReporitoriesImpl{Db: db}

}

func (u *UserReporitoriesImpl) CreateUser(user domain.User) error {
	if err := u.Db.Create(&user).Error; err != nil {
		log.Printf("failed to create user with email %s: %v", user.Email, err)
		return err
	}
	log.Printf("User %s successfully registered", user.Email)
	return nil
}

func (u *UserReporitoriesImpl) FindById(userId int) (*domain.User, error) {
	var user domain.User

	result := u.Db.Find(&user, userId)

	if result == nil {
		return nil, errors.New("User Not Found within ID")

	} else {
		return &user, nil
	}
}

func (u *UserReporitoriesImpl) FindOrCreateByEmail(email string, googleID string) (*domain.User, error) {
	var user domain.User

	result := u.Db.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user = domain.User{
			Email: email,
		}

		if googleID != "" {
			user.GoogleID = googleID
		}

		if err := u.Db.Create(&user).Error; err != nil {
			return nil, err
		}
	} else if result.Error != nil {

		return nil, result.Error
	}

	return &user, nil
}

func (u *UserReporitoriesImpl) CreateSession(session domain.Session) error {
	u.Db.Where("user_id", session.UserID).Delete(&domain.Session{})

	if err := u.Db.Create(&session).Error; err != nil {
		log.Printf("Failed to create session for user %d: %v", session.UserID, err)
		return err
	}

	log.Printf("Session created successfully for user %d", session.UserID)
	return nil
}

func (u *UserReporitoriesImpl) FindSessionByUserID(userID uint) (*domain.Session, error) {
	var session domain.Session

	result := u.Db.Where("user_id", userID).First(&session)

	if result.Error != nil {

		return nil, errors.New("session expired or not found")
	}

	return &session, nil
}
