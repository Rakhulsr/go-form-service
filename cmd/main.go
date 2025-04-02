package main

import (
	"log"
	"net/http"

	"github.com/Rakhulsr/go-form-service/db"
	"github.com/Rakhulsr/go-form-service/internal/routes"
)

func main() {
	db, err := db.NewDbConnection()
	if err != nil {
		log.Fatal(err)
	}

	router := routes.NewRouter(db)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server up at port 8080")

	if err := server.ListenAndServe(); err != nil {
		log.Println("failed to connecting to the server")
	}
}
