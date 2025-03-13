package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var db *sql.DB

func generateUniqueNumber() int64 {
	return time.Now().UnixNano()
}

func connectDB() (*sql.DB, error) {
	dsn := "user=username password=password dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func checkEntry(db *sql.DB, table string, column string, value string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE %s = $1 LIMIT 1", table, column)
	var exists int
	row := db.QueryRow(query, value)
	err := row.Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func toBase62(num int64) string {
	if num == 0 {
		return "0"
	}
	base := int64(len(base62Chars))
	result := ""
	for num > 0 {
		remainder := num % base
		result = string(base62Chars[remainder]) + result
		num /= base
	}
	return result
}

func generateBase62ID() string {
	uniqueNum := generateUniqueNumber()
	return toBase62(uniqueNum)
}

func genRandNum(db *sql.DB, table string, column string) (string, error) {
	num := generateBase62ID()
	res, err := checkEntry(db, table, column, num)
	if err != nil {
		return "", err
	}
	if res {
		return genRandNum(db, table, column)
	}
	return num, nil
}

func Shorten(w http.ResponseWriter, r *http.Request) {
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
	shortCode, err := genRandNum(db, "urls", "short_code")
	if err != nil {
		http.Error(w, "Failed to generate short code", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", shortCode, longUrl)
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

func Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]
	var longUrl string
	query := db.QueryRow("SELECT original_url FROM urls WHERE short_code = $1", shortCode)
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

func main() {
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/shorten", Shorten).Methods("POST")
	r.HandleFunc("/{shortCode}", Redirect).Methods("GET")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
