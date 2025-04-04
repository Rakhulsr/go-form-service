package utils

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Rakhulsr/go-form-service/config/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JwtSecret = []byte(env.ENV.Secret)

type Claims struct {
	Email    string `json:"email"`
	GoogleID string `json:"google_id"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccesToken(email, googleID string, userID uint) (string, error) {
	claims := Claims{
		Email:    email,
		GoogleID: googleID,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func GenerateRefreshToken(email, googleID string, userID uint) (string, error) {
	claims := Claims{
		Email:    email,
		GoogleID: googleID,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func ExtractUserIDFromExpiredToken(r *http.Request) (uint, error) {

	tokenCookie, err := r.Cookie("access-token")
	if err != nil {
		return 0, errors.New("token tidak ditemukan")
	}

	tokenString := tokenCookie.Value
	log.Println(tokenString)

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return 0, errors.New("gagal memparsing token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("gagal mendapatkan claims token")
	}

	var userIDFloat float64
	if val, exists := claims["userID"]; exists {
		userIDFloat = val.(float64)
	} else if val, exists := claims["user_id"]; exists {
		userIDFloat = val.(float64)
	} else {
		return 0, errors.New("user ID tidak ditemukan dalam token")
	}

	userID := uint(userIDFloat)

	return userID, nil
}
