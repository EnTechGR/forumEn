package utils

import (
	"errors"
	"regexp"
)

func ValidateEmail(email string) error {
	// Email validation
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("Invalid email format")
	}

	return nil
}
