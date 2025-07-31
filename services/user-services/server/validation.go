package server

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/ttacon/libphonenumber"
)

// MARK: Validation functions for user input
func ValidateFullName(fullName string) bool {
	if len(strings.TrimSpace(fullName)) < 2 {
		return false
	}
	if len(fullName) > 100 {
		return false
	}
	for _, char := range fullName {
		if !unicode.IsLetter(char) && char != ' ' && char != '-' && char != '\'' {
			return false
		}
	}
	hasLetter := false
	for _, char := range fullName {
		if unicode.IsLetter(char) {
			hasLetter = true
			break
		}
	}
	return hasLetter
}

// MARK: Validation functions for user input
func ValidateEmail(email string) bool {
	if len(email) < 5 || len(email) > 254 {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// MARK: Validation functions for user input
func ValidatePhoneNumber(phoneNumber string) bool {
	if len(phoneNumber) == 0 {
		return true
	}
	parsedNumber, err := libphonenumber.Parse(phoneNumber, "")
	if err != nil {
		return false
	}
	return libphonenumber.IsValidNumber(parsedNumber)
}

// MARK: Validation functions for user input
func ValidatePhoneNumberWithCountry(phoneNumber, countryCode string) bool {
	if len(phoneNumber) == 0 {
		return true
	}
	if len(countryCode) == 0 {
		return ValidatePhoneNumber(phoneNumber)
	}

	var fullPhoneNumber string
	var regionCode string

	if strings.HasPrefix(countryCode, "+") {

		cleanPhoneNumber := strings.TrimPrefix(phoneNumber, "0")
		fullPhoneNumber = countryCode + cleanPhoneNumber
		regionCode = ""
	} else if len(countryCode) > 2 {
		if matched, _ := regexp.MatchString(`^\d+$`, countryCode); matched {
			regionCode = ""
			fullPhoneNumber = "+" + countryCode + phoneNumber
		} else {
			regionCode = countryCode
			fullPhoneNumber = phoneNumber
		}
	} else {
		regionCode = countryCode
		fullPhoneNumber = phoneNumber
	}

	parsedNumber, err := libphonenumber.Parse(fullPhoneNumber, regionCode)
	if err != nil {
		return false
	}
	return libphonenumber.IsValidNumber(parsedNumber)
}

// MARK: Format phone number with country code
func FormatPhoneNumber(phoneNumber, countryCode string) (string, error) {
	regionCode := countryCode
	if strings.HasPrefix(countryCode, "+") {
		regionCode = ""
	} else if len(countryCode) > 2 {
		if matched, _ := regexp.MatchString(`^\d+$`, countryCode); matched {
			regionCode = ""
			phoneNumber = "+" + countryCode + phoneNumber
		}
	}

	parsedNumber, err := libphonenumber.Parse(phoneNumber, regionCode)
	if err != nil {
		return "", err
	}

	return libphonenumber.Format(parsedNumber, libphonenumber.INTERNATIONAL), nil
}

// MARK: Validate country code
func ValidateCountryCode(countryCode string) bool {
	if len(countryCode) == 0 {
		return true // Country code is optional
	}
	if len(countryCode) == 2 {
		countryRegex := regexp.MustCompile(`^[A-Z]{2}$`)
		return countryRegex.MatchString(strings.ToUpper(countryCode))
	}
	countryRegex := regexp.MustCompile(`^\+?[1-9]\d{0,3}$`)
	return countryRegex.MatchString(countryCode)
}

// MARK: Validate password
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	if len(password) > 128 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}
