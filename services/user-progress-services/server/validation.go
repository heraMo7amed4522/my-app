package server

import (
	"regexp"
	"strings"
)

// ValidateUserID validates user ID format
func ValidateUserID(userID string) bool {
	userID = strings.TrimSpace(userID)
	return len(userID) > 0 && len(userID) <= 100
}

// ValidateTemplateID validates template ID format
func ValidateTemplateID(templateID string) bool {
	templateID = strings.TrimSpace(templateID)
	return len(templateID) > 0 && len(templateID) <= 100
}

// ValidateSectionID validates section ID format
func ValidateSectionID(sectionID string) bool {
	sectionID = strings.TrimSpace(sectionID)
	return len(sectionID) > 0 && len(sectionID) <= 100
}

// ValidateProgress validates progress value
func ValidateProgress(progress float32) bool {
	return progress >= 0 && progress <= 100
}

// ValidateUUID validates UUID format
func ValidateUUID(uuid string) bool {
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidPattern.MatchString(uuid)
}

// ValidatePagination validates pagination parameters
func ValidatePagination(page, limit int32) (int32, int32) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	return page, limit
}