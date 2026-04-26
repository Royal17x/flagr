package validator

import (
	"fmt"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if len(email) > 254 {
		return fmt.Errorf("email too long")
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 128 {
		return fmt.Errorf("password too long")
	}
	return nil
}

func ValidateOrgName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("org_name is required")
	}
	if len(name) > 100 {
		return fmt.Errorf("org_name too long")
	}
	return nil
}

func ValidateFlagKey(key string) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}
	validKey := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	if !validKey.MatchString(key) {
		return fmt.Errorf("key must contain only letters, numbers, hyphens and underscores")
	}
	if len(key) > 100 {
		return fmt.Errorf("key too long")
	}
	return nil
}
