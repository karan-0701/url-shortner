package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

// GetIPAddress extracts the user's IP address from the HTTP request
func GetIPAddress(r *http.Request) string {
	// Try X-Forwarded-For header first (common for clients behind proxy)
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		ip = strings.TrimSpace(ips[0]) // Get the client's original IP
		return ip
	}

	// Try X-Real-IP header (used by some proxies)
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fall back to RemoteAddr
	ip = r.RemoteAddr

	// Handle IPv6 addresses which have format [IPv6]:port
	if strings.HasPrefix(ip, "[") {
		// Extract IPv6 address without port
		end := strings.LastIndex(ip, "]")
		if end > 0 {
			return ip[1:end]
		}
	} else if strings.Contains(ip, ":") {
		// Handle IPv4 addresses which have format IPv4:port
		return strings.Split(ip, ":")[0]
	}

	return ip
}

// GetGeolocation fetches the country and city for a given IP address using the ipstack API
func GetGeolocation(ip string, apiKey string) (string, string, error) {
	url := fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s", ip, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("failed to call ipstack API: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		CountryName string `json:"country_name"`
		City        string `json:"city"`
		Success     bool   `json:"success"`
		Error       struct {
			Info string `json:"info"`
		} `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode ipstack response: %w", err)
	}

	if !result.Success {
		return "", "", fmt.Errorf("ipstack API error: %s", result.Error.Info)
	}

	return result.CountryName, result.City, nil
}
