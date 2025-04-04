package home

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Rakhulsr/go-form-service/utils"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		log.Println("Failed to get cookie in home handler")
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	claims, err := utils.VerifyToken(cookie.Value)
	if err != nil {
		log.Println("Failed to verify token in home handler:", err)
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("../templates/home/home.html")
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	data := struct {
		Email    string
		GoogleID string
	}{
		Email:    claims.Email,
		GoogleID: claims.GoogleID,
	}

	tmpl.Execute(w, data)
}
