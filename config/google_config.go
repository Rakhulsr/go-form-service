package config

import (
	"github.com/Rakhulsr/go-form-service/config/env"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleConfig struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig GoogleConfig

func GoogleOathConfig() oauth2.Config {
	AppConfig.GoogleLoginConfig = oauth2.Config{
		ClientID:     env.ENV.GoogleCliendID,
		ClientSecret: env.ENV.GoogleClientSecret,
		RedirectURL:  env.ENV.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	return AppConfig.GoogleLoginConfig

}
