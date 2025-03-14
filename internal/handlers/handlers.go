package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

	_, err = h.DB.Exec("INSERT INTO urls (short_code, original_url) VALUES (?, ?)", shortCode, longUrl)
	if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"short_url": fmt.Sprintf("http://localhost:8080/%s", shortCode),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Redirect handles the redirection request
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	var longUrl string
	query := h.DB.QueryRow("SELECT original_url FROM urls WHERE short_code = ?", shortCode)
	err := query.Scan(&longUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "URL not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, longUrl, http.StatusMovedPermanently)
}
