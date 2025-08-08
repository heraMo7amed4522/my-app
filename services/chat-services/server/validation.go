package server

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// ValidateMessageContent validates message content
func ValidateMessageContent(content string) bool {
	if len(strings.TrimSpace(content)) == 0 {
		return false
	}
	if utf8.RuneCountInString(content) > 4000 { // Max 4000 characters
		return false
	}
	return true
}

// ValidateUserID validates user ID format
func ValidateUserID(userID string) bool {
	if len(strings.TrimSpace(userID)) == 0 {
		return false
	}
	// UUID format validation
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(userID)
}

// ValidateGroupName validates group name
func ValidateGroupName(name string) bool {
	name = strings.TrimSpace(name)
	if len(name) < 1 || len(name) > 100 {
		return false
	}
	return true
}

// ValidateGroupDescription validates group description
func ValidateGroupDescription(description string) bool {
	if utf8.RuneCountInString(description) > 500 {
		return false
	}
	return true
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	if len(email) < 5 || len(email) > 254 {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateLimit validates pagination limit
func ValidateLimit(limit int32) int32 {
	if limit <= 0 || limit > 100 {
		return 20 // Default limit
	}
	return limit
}

// ValidateOffset validates pagination offset
func ValidateOffset(offset int32) int32 {
	if offset < 0 {
		return 0
	}
	return offset
}

// ValidateSearchQuery validates search query
func ValidateSearchQuery(query string) bool {
	query = strings.TrimSpace(query)
	if len(query) < 1 || len(query) > 100 {
		return false
	}
	return true
}

// ValidatePollQuestion validates poll question
func ValidatePollQuestion(question string) bool {
	question = strings.TrimSpace(question)
	if len(question) < 1 || len(question) > 500 {
		return false
	}
	return true
}

// ValidatePollOptions validates poll options
func ValidatePollOptions(options []string) bool {
	if len(options) < 2 || len(options) > 10 {
		return false
	}
	for _, option := range options {
		option = strings.TrimSpace(option)
		if len(option) < 1 || len(option) > 100 {
			return false
		}
	}
	return true
}
