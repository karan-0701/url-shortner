package main

import (
	"database/sql"
	"fmt"
	"time"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func generateUniqueNumber() int64 {
	return time.Now().UnixNano()
}
func connectDB() (*sql.DB, error) {
	dsn := "user=username password=password dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func checkEntry(db *sql.DB, table string, column string, value string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WEHERE %s = $1 LIMIT 1", table, column)
	var exists int
	row := db.QueryRow(query, value)
	// store 1 in exts if there exist a row in the database
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
		return " ", nil
	}

	if res == true {
		return genRandNum(db, table, column)
	}

	return num, nil
}

func main() {
	r := mux.NewRouter()
}
