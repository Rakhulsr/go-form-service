package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/Rakhulsr/go-form-service/internal/services"
	"github.com/Rakhulsr/go-form-service/utils"
)

type Middleware struct {
	UserService services.UserService
}

func NewMiddlewareImpl(userService services.UserService) *Middleware {
	return &Middleware{UserService: userService}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		accessCookie, err := r.Cookie("access-token")
		if err == nil {
			claims, err := utils.VerifyToken(accessCookie.Value)
			if err == nil {
				log.Println(claims)
				next.ServeHTTP(w, r)
				return
			}
		}

		refreshCookie, err := r.Cookie("refresh-token")
		if err == nil {
			claims, err := utils.VerifyToken(refreshCookie.Value)
			if err == nil {
				newAccessToken, _ := utils.GenerateAccesToken(claims.Email, claims.GoogleID, claims.UserID)
				http.SetCookie(w, &http.Cookie{
					Name:     "access-token",
					Value:    newAccessToken,
					Expires:  time.Now().Add(15 * time.Minute),
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteStrictMode,
				})
				next.ServeHTTP(w, r)
				return
			}
		}

		m.RefreshAccessToken(w, r, next)
	})
}

func (m *Middleware) RefreshAccessToken(w http.ResponseWriter, r *http.Request, next http.Handler) {

	userID, err := utils.ExtractUserIDFromExpiredToken(r)
	if err != nil {

		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}

	session, err := m.UserService.RefreshSession(userID)
	if err != nil || session == nil {

		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}

	claims, err := utils.VerifyToken(session.Token)
	if err != nil {

		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}

	newAccessToken, err := utils.GenerateAccesToken(claims.Email, claims.GoogleID, claims.UserID)
	if err != nil {

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    newAccessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	next.ServeHTTP(w, r)
}
