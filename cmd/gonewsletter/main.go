package main

import (
	"log"
	"net/http"

	"github.com/edisss1/go-newsletter/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	http.HandleFunc("/", server.FormHandler)
	http.HandleFunc("/submit", server.SubmitHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Couldn't start server: %v", err)
	}
}
