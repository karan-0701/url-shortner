package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	// PostgreSQL connection string
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fallback to a local connection string for development
		connStr = "user=username dbname=mydb sslmode=disable password=password"
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create the urls table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_code TEXT NOT NULL UNIQUE,
			original_url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Create url_analytics table if it doesn't exist
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
