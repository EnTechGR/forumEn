package repository

import (
	"database/sql"
	"time"

	"forum/config"
	"forum/models"
	"forum/utils"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create adds a new user to the database
func (r *UserRepository) Create(reg models.UserRegistration) (*models.User, error) {
	// Check if email is already taken
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", reg.Email).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, config.ErrEmailTaken
	}

	// Check if username is already taken
	err = r.DB.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", reg.Username).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, config.ErrUsernameTaken
	}

	// Start a transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Generate UUID for the user
	userID := utils.GenerateUUID()

	// Set creation time once and reuse it
	createdAt := time.Now()

	// Insert user record
	_, err = tx.Exec(
		"INSERT INTO user (user_id, username, email, created_at) VALUES (?, ?, ?, ?)",
		userID, reg.Username, reg.Email, createdAt,
	)
	if err != nil {
		return nil, err
	}

	// Hash the password
	passwordHash, err := utils.HashPassword(reg.Password)
	if err != nil {
		return nil, err
	}

	// Insert authentication record
	_, err = tx.Exec(
		"INSERT INTO user_auth (user_id, password_hash) VALUES (?, ?)",
		userID, passwordHash,
	)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Return the newly created user
	user := &models.User{
		ID:        userID,
		Username:  reg.Username,
		Email:     reg.Email,
		CreatedAt: createdAt,
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.DB.QueryRow(
		"SELECT user_id, username, email, created_at FROM user WHERE LOWER(email) = LOWER(?)",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, config.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	var user models.User

	err := r.DB.QueryRow(
		"SELECT user_id, username, email, created_at FROM user WHERE user_id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, config.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetAuthByUserID retrieves user authentication data by user ID
func (r *UserRepository) GetAuthByUserID(userID string) (*models.UserAuth, error) {
	var auth models.UserAuth

	err := r.DB.QueryRow(
		"SELECT user_id, password_hash FROM user_auth WHERE user_id = ?",
		userID,
	).Scan(&auth.UserID, &auth.PasswordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, config.ErrUserNotFound
		}
		return nil, err
	}

	return &auth, nil
}

// Authenticate validates a user's login credentials
func (r *UserRepository) Authenticate(login models.UserLogin) (*models.User, error) {
	// Get the user by email
	user, err := r.GetByEmail(login.Email)
	if err != nil {
		return nil, config.ErrInvalidCredentials
	}

	// Get the user's authentication data
	auth, err := r.GetAuthByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	// Check the password
	if !utils.CheckPasswordHash(login.Password, auth.PasswordHash) {
		return nil, config.ErrInvalidCredentials
	}

	return user, nil
}
