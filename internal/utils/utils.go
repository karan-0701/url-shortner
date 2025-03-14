package utils

import (
	"database/sql"
	"fmt"
	"time"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateUniqueNumber generates a unique number based on the current time
func GenerateUniqueNumber() int64 {
	return time.Now().UnixNano()
}

// ToBase62 converts a number to a Base62 string
func ToBase62(num int64) string {
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

// GenerateBase62ID generates a unique Base62 ID
func GenerateBase62ID() string {
	uniqueNum := GenerateUniqueNumber()
	return ToBase62(uniqueNum)
}

// CheckEntry checks if an entry exists in the database
func CheckEntry(db *sql.DB, table string, column string, value string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE %s = ? LIMIT 1", table, column)
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

// GenRandNum generates a random unique Base62 ID
func GenRandNum(db *sql.DB, table string, column string) (string, error) {
	num := GenerateBase62ID()
	res, err := CheckEntry(db, table, column, num)
	if err != nil {
		return "", err
	}
	if res {
		return GenRandNum(db, table, column)
	}
	return num, nil
}
