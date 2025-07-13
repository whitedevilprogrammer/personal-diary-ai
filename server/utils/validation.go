package utils

import (
	"crypto/rand"
	"encoding/base64"
	"regexp"
	"strings"
)

// GenerateRandomString generates a cryptographically secure random string
// of the specified length using base64 encoding
func GenerateRandomString(length int) string {
	// Calculate the number of bytes needed
	// Base64 encoding uses 4 characters for every 3 bytes
	byteLength := (length * 3) / 4
	if (length*3)%4 != 0 {
		byteLength++
	}
	
	// Generate random bytes
	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to a simpler method if crypto/rand fails
		return generateFallbackString(length)
	}
	
	// Encode to base64 and trim to desired length
	encoded := base64.URLEncoding.EncodeToString(bytes)
	encoded = strings.ReplaceAll(encoded, "=", "") // Remove padding
	
	if len(encoded) > length {
		return encoded[:length]
	}
	
	return encoded
}

// generateFallbackString provides a fallback random string generation
func generateFallbackString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	
	// Use crypto/rand even for fallback
	randBytes := make([]byte, length)
	rand.Read(randBytes)
	
	for i, b := range randBytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	
	return string(bytes)
}

// IsValidEmail validates email format using regex
func IsValidEmail(email string) bool {
	// Basic email validation regex
	// This covers most common email formats
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	
	// Additional checks
	if len(email) == 0 || len(email) > 254 {
		return false
	}
	
	// Check for valid format
	if !emailRegex.MatchString(email) {
		return false
	}
	
	// Additional validations
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	localPart := parts[0]
	domainPart := parts[1]
	
	// Local part validations
	if len(localPart) == 0 || len(localPart) > 64 {
		return false
	}
	
	// Domain part validations
	if len(domainPart) == 0 || len(domainPart) > 253 {
		return false
	}
	
	// Check for consecutive dots
	if strings.Contains(email, "..") {
		return false
	}
	
	// Check if starts or ends with dot
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return false
	}
	
	return true
}

// ValidatePassword validates password strength
func ValidatePassword(password string) (bool, string) {
	if len(password) < 6 {
		return false, "Password must be at least 6 characters long"
	}
	
	if len(password) > 128 {
		return false, "Password must not exceed 128 characters"
	}
	
	// Check for at least one letter and one number (optional, uncomment if needed)
	/*
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	
	if !hasLetter {
		return false, "Password must contain at least one letter"
	}
	
	if !hasNumber {
		return false, "Password must contain at least one number"
	}
	*/
	
	return true, ""
}

// SanitizeString removes potentially harmful characters from strings
func SanitizeString(input string) string {
	// Remove null bytes and control characters
	sanitized := strings.ReplaceAll(input, "\x00", "")
	
	// Remove other control characters except tab, newline, and carriage return
	result := ""
	for _, char := range sanitized {
		if char >= 32 || char == 9 || char == 10 || char == 13 {
			result += string(char)
		}
	}
	
	return strings.TrimSpace(result)
}

// GenerateSecureToken generates a secure token for various purposes
func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	
	return base64.URLEncoding.EncodeToString(bytes), nil
}