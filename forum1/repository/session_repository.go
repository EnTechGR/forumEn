package repository

import (
	"database/sql"
	"time"

	"forum/config"
	"forum/models"
	"forum/utils"
)

// SessionRepository handles session-related database operations
type SessionRepository struct {
	DB *sql.DB
}

// NewSessionRepository creates a new SessionRepository
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

// Create creates a new session for a user
func (r *SessionRepository) Create(userID, ipAddress string) (*models.Session, error) {
	// Start a transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// First, delete any existing sessions for this user
	_, err = tx.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	// Generate a new session ID
	sessionID := utils.GenerateSessionToken()
	expiresAt := utils.CalculateSessionExpiry()
	now := time.Now()

	// Insert the new session
	_, err = tx.Exec(
		"INSERT INTO sessions (user_id, session_id, ip_address, created_at, expires_at) VALUES (?, ?, ?, ?, ?)",
		userID, sessionID, ipAddress, now, expiresAt,
	)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Return the session
	session := &models.Session{
		UserID:    userID,
		SessionID: sessionID,
		IPAddress: ipAddress,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}

	return session, nil
}

// GetBySessionID retrieves a session by its ID
func (r *SessionRepository) GetBySessionID(sessionID string) (*models.Session, error) {

	var session models.Session

	err := r.DB.QueryRow(
		"SELECT user_id, session_id, ip_address, created_at, expires_at FROM sessions WHERE session_id = ?",
		sessionID,
	).Scan(&session.UserID, &session.SessionID, &session.IPAddress, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, config.ErrSessionNotFound
		}
		return nil, err
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// Delete the expired session
		_, _ = r.DB.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
		return nil, config.ErrSessionExpired
	}

	return &session, nil
}

// Delete removes a session
func (r *SessionRepository) Delete(sessionID string) error {
	_, err := r.DB.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}
