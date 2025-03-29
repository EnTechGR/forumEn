package utils

import (
	"errors"
	"forum/config"
	"regexp"
)

func ValidatePassword(password string) error {
	// Password validation
	if len(password) < config.MIN_PASSWORD_LEN || len(password) > config.MAX_PASSWORD_LEN {
		return errors.New("Password must be at least 8 characters long")
	}

	// Check password complexity
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	if !(hasUppercase && hasLowercase && hasNumber && hasSpecial) {
		return errors.New("Password must contain uppercase, lowercase, number, and special character")
	}

	return nil
}
