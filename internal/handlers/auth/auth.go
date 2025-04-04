package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/Rakhulsr/go-form-service/config"
	"github.com/Rakhulsr/go-form-service/internal/models/domain"
	"github.com/Rakhulsr/go-form-service/internal/services"
	"golang.org/x/oauth2"
)

type OauthHandlerImpl struct {
	UserService services.UserService
}

func NewOAuthHandler(userService services.UserService) *OauthHandlerImpl {
	return &OauthHandlerImpl{
		UserService: userService,
	}
}

func (h *OauthHandlerImpl) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthConfig := config.GoogleOathConfig()

	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)

}

func (h *OauthHandlerImpl) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	oauthConfig := config.GoogleOathConfig()
	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(w, "Missing state parameter", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	log.Println("Access Token:", token.AccessToken)

	userData, err := getGoogleUserData(token.AccessToken)
	if err != nil {
		log.Printf("Failed to get user data: %v", err)
		http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		return
	}

	log.Println("User Data:", userData)
	fmt.Println(userData.Email, userData.GoogleID)

	exUser, err := h.UserService.FindOrCreateByEmail(userData.Email, userData.GoogleID)
	if err != nil {
		log.Printf("Error finding or creating user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := h.UserService.LoginGoogle(userData.Email, userData.GoogleID, exUser.ID)

	if err != nil {
		log.Printf("Error generating tokens: %v", err)
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   7 * 24 * 60 * 60,
	})
	w.Header().Set("Authorization", "Bearer "+accessToken)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func getGoogleUserData(accessToken string) (*domain.User, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	googleID, ok := responseMap["id"].(string)
	if !ok || googleID == "" {
		return nil, fmt.Errorf("Google ID (id) not found in response")
	}

	email, ok := responseMap["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("email not found in response")
	}

	userData := &domain.User{
		GoogleID: googleID,
		Email:    email,
	}

	return userData, nil
}

func (h *OauthHandlerImpl) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../templates/auth/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (h *OauthHandlerImpl) Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}
