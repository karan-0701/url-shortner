package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/karan-0701/url-shortner/internal/utils"
)

// Handler struct holds the database connection
type Handler struct {
	DB *sql.DB
}

// NewHandler returns a new Handler instance
func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

// Shorten handles the URL shortening request
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	longUrl := r.FormValue("url")
	if longUrl == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortCode, err := utils.GenRandNum(h.DB, "urls", "short_code")
	if err != nil {
		http.Error(w, "Failed to generate short code", http.StatusInternalServerError)
		return
	}

	// Update the query to use PostgreSQL placeholders
	_, err = h.DB.Exec("INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", shortCode, longUrl)
	if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	baseURL := fmt.Sprintf("https://%s", r.Host)

	response := map[string]string{
		"short_url": fmt.Sprintf("%s/%s", baseURL, shortCode),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Redirect handles the redirection request
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	var longUrl string
	// Update the query to use PostgreSQL placeholders
	query := h.DB.QueryRow("SELECT original_url FROM urls WHERE short_code = $1", shortCode)
	err := query.Scan(&longUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "URL not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	referrer := r.Referer()
	if referrer == "" {
		referrer = "direct"
	}

	// Call ipstack API to get geolocation
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Update the query to use PostgreSQL placeholders
	_, err = h.DB.Exec("INSERT INTO url_analytics (short_code, referrer) VALUES ($1, $2)", shortCode, referrer)
	if err != nil {
		log.Println("Failed to log click:", err)
	}

	http.Redirect(w, r, longUrl, http.StatusMovedPermanently)
}
