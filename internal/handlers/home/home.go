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

	token, err := utils.VerifyToken(cookie.Value)
	if err != nil || !token.Valid {
		log.Println("Failed to veirfy in home handler")

		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	claims, ok := token.Claims.(*utils.Claims)
	if !ok {
		log.Println("Failed to claims in home handler")

		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	email := claims.Email
	googleID := claims.GoogleID

	tmpl, err := template.ParseFiles("../templates/home/home.html")
	if err != nil {
		http.Error(w, "Gagal memuat halaman", http.StatusInternalServerError)
		return
	}

	data := struct {
		Email    string
		GoogleID string
	}{
		Email:    email,
		GoogleID: googleID,
	}

	tmpl.Execute(w, data)
}
