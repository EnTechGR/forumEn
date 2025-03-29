package config

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrEmailTaken           = errors.New("email is already taken")
	ErrUsernameTaken        = errors.New("username is already taken")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionExpired       = errors.New("session expired")
)
