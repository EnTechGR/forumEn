package models

import "time"

// User represents a forum user
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}


// UserAuth contains user authentication information
type UserAuth struct {
	UserID       string `json:"-"`
	PasswordHash string `json:"-"`
}

// UserLogin is used for login requests
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
