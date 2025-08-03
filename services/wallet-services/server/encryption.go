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
	"strconv"
	"strings"
)

// getEncryptionKey derives a 32-byte key from environment variable or default
func getEncryptionKey() []byte {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		// Use a default key for development - in production, this should be from env
		key = "wallet-service-encryption-key-2024"
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

// DecryptData decrypts data encrypted with EncryptData
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

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}

	return string(plaintext), nil
}

// MaskBalance masks wallet balance for logging/display purposes
func MaskBalance(balance string) string {
	if len(balance) <= 4 {
		return strings.Repeat("*", len(balance))
	}
	return strings.Repeat("*", len(balance)-4) + balance[len(balance)-4:]
}

// ValidateAmount validates if the amount string is a valid positive decimal
func ValidateAmount(amount string) bool {
	if amount == "" {
		return false
	}

	// Parse as float to validate format
	val, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return false
	}

	// Amount must be positive
	return val > 0
}

// ValidateCurrency validates if the currency code is supported
func ValidateCurrency(currency string) bool {
	supportedCurrencies := map[string]bool{
		"USD": true,
		"EUR": true,
		"GBP": true,
		"JPY": true,
		"CAD": true,
		"AUD": true,
	}

	return supportedCurrencies[strings.ToUpper(currency)]
}