package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

// getEncryptionKey derives a 32-byte key from environment variable or default
func getEncryptionKey() []byte {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		// Use a default key for development - in production, this should be from env
		key = "card-service-encryption-key-2024"
	}
	
	// Create a 32-byte key using SHA256
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

// EncryptData encrypts sensitive data using AES-GCM
func EncryptData(plaintext string) (string, error) {
	if plaintext == "" {
		return "", errors.New("plaintext cannot be empty")
	}

	key := getEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptData decrypts data that was encrypted with EncryptData
func DecryptData(encryptedData string) (string, error) {
	if encryptedData == "" {
		return "", errors.New("encrypted data cannot be empty")
	}

	key := getEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}

	return string(plaintext), nil
}

// MaskCardNumber masks a card number showing only last 4 digits
func MaskCardNumber(cardNumber string) string {
	if len(cardNumber) < 4 {
		return "****"
	}
	lastFour := cardNumber[len(cardNumber)-4:]
	return "****-****-****-" + lastFour
}

// ValidateCardNumber performs basic card number validation
func ValidateCardNumber(cardNumber string) bool {
	// Remove spaces and dashes
	cleaned := ""
	for _, char := range cardNumber {
		if char >= '0' && char <= '9' {
			cleaned += string(char)
		}
	}

	// Check length (13-19 digits for most cards)
	if len(cleaned) < 13 || len(cleaned) > 19 {
		return false
	}

	// Luhn algorithm check
	return luhnCheck(cleaned)
}

// luhnCheck implements the Luhn algorithm for card validation
func luhnCheck(cardNumber string) bool {
	sum := 0
	alternate := false

	// Process digits from right to left
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}
		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}