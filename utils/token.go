package utils

import (
	"errors"
	"time"

	"github.com/Rakhulsr/go-form-service/config/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JwtSecret = []byte(env.ENV.SecretJWT)

type Claims struct {
	Email    string `json:"email"`
	GoogleID string `json:"google_id"`
	jwt.RegisteredClaims
}

func GenerateAccesToken(email, googleID string) (string, error) {
	claims := Claims{
		Email:    email,
		GoogleID: googleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func GenerateRefreshToken(email, googleID string) (string, error) {
	claims := Claims{
		Email:    email,
		GoogleID: googleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(*Claims); ok && token.Valid {
		return token, nil
	}

	return nil, errors.New("invalid token claims")
}
