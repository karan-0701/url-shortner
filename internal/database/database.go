package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() (*sql.DB, error) {
	// Connect to SQLite database (creates the file if it doesn't exist)
	db, err := sql.Open("sqlite3", "./urlshortener.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create the urls table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			short_code TEXT NOT NULL UNIQUE,
			original_url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Crete urls analytics if it does not exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS url_analytics (
			id SERIAL PRIMARY KEY,
			short_code TEXT NOT NULL,
			referrer TEXT,
			country TEXT,
			city TEXT,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (short_code) REFERENCES urls(short_code)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Connected to the database successfully")
	return db, nil
}
