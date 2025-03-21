package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/karan-0701/url-shortner/internal/database"
	"github.com/karan-0701/url-shortner/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Initialize the database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	handler := handlers.NewHandler(db)

	// Set up routes
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("frontend"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "frontend/index.html") }).Methods("GET")
	r.HandleFunc("/shorten", handler.Shorten).Methods("POST")
	r.HandleFunc("/{shortCode}", handler.Redirect).Methods("GET")

	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
