package utils

import (
	"errors"
	"forum/config"
	"regexp"
)

func ValidateUsername(username string) error {
	// Username validation
	if len(username) < config.MIN_USERNAME_LEN || len(username) > config.MAX_USERNAME_LEN {
		return errors.New("Username must be between 3 and 30 characters long")

	}

	// Username character validation (alphanumeric and underscore only)
	usernameRegex := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	if !usernameRegex.MatchString(username) {
		return errors.New("Username can only contain alphanumeric characters and underscores")
	}

	return nil
}
