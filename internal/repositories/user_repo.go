package repositories

import (
	"errors"
	"log"

	"github.com/Rakhulsr/go-form-service/internal/models/domain"
	"gorm.io/gorm"
)

type UserRepositories interface {
	FindById(userId int) (*domain.User, error)
	// FindByEmail(email string) (*domain.User, error)
	CreateUser(user domain.User) error
	FindOrCreateByEmail(email string, googleID string) (*domain.User, error)
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

	// Coba cari user berdasarkan email
	result := u.Db.Where("email = ?", email).First(&user)

	// Jika user tidak ditemukan, buat user baru
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user = domain.User{
			Email: email,
		}

		// Masukkan googleID hanya jika tidak kosong
		if googleID != "" {
			user.GoogleID = googleID
		}

		if err := u.Db.Create(&user).Error; err != nil {
			return nil, err
		}
	} else if result.Error != nil {
		// Jika terjadi error selain "record not found", kembalikan error
		return nil, result.Error
	}

	return &user, nil
}
