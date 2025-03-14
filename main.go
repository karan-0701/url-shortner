package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/karan-0701/url-shortner/internal/database"
	"github.com/karan-0701/url-shortner/internal/handlers"
)

func main() {
	// Initialize the database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	handler := handlers.NewHandler(db)

	// Set up routes
	r := mux.NewRouter()
	r.HandleFunc("/shorten", handler.Shorten).Methods("POST")
	r.HandleFunc("/{shortCode}", handler.Redirect).Methods("GET")

	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
