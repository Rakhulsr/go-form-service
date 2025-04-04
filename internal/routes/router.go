package routes

import (
	"github.com/Rakhulsr/go-form-service/internal/handlers/auth"
	"github.com/Rakhulsr/go-form-service/internal/handlers/home"
	"github.com/Rakhulsr/go-form-service/internal/repositories"
	"github.com/Rakhulsr/go-form-service/internal/services"
	"github.com/Rakhulsr/go-form-service/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	validator := validator.New()

	//auth
	oauthRouter := router.PathPrefix("/auth").Subrouter()
	oauthRepo := repositories.NewUserRepositoryImpl(db)
	oauthService := services.NewUserServiceImpl(&oauthRepo, validator)
	oauthHandler := auth.NewOAuthHandler(oauthService)

	oauthRouter.HandleFunc("/login", oauthHandler.LoginPage).Methods("GET")
	oauthRouter.HandleFunc("/google/login", oauthHandler.GoogleLogin).Methods("GET")
	oauthRouter.HandleFunc("/google/callback", oauthHandler.GoogleCallback).Methods("GET")
	oauthRouter.HandleFunc("/logout", oauthHandler.Logout).Methods("POST", "GET")

	//home
	middleware := middlewares.NewMiddlewareImpl(oauthService)
	homeRouter := router.PathPrefix("/").Subrouter()
	homeRouter.Use(middleware.AuthMiddleware)

	homeRouter.HandleFunc("/home", home.HomePage).Methods("GET")

	return router

}
