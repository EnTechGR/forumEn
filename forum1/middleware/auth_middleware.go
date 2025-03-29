package middleware

import (
	"context"
	"log"
	"net/http"

	"forum/models"
	"forum/repository"
)

// Authentication middleware checks if the user is authenticated
type AuthMiddleware struct {
	SessionRepo *repository.SessionRepository
	UserRepo    *repository.UserRepository
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(sessionRepo *repository.SessionRepository, userRepo *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
	}
}

// Authenticate middleware verifies authentication and sets user in context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the session cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Check if session exists in database directly
		var count int
		err = m.SessionRepo.DB.QueryRow("SELECT COUNT(*) FROM sessions WHERE session_id = ?", cookie.Value).Scan(&count)
		if err != nil {
			log.Printf("Error checking session in DB: %v", err)
		} else {
			log.Printf("Session in database: %v (count=%d)", count > 0, count)
		}

		// Validate the session
		session, err := m.SessionRepo.GetBySessionID(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Get the user
		user, err := m.UserRepo.GetByID(session.UserID)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		log.Printf("User retrieved successfully: %+v", user)

		// Set user in context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth middleware ensures the user is authenticated
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetCurrentUser returns the authenticated user from the context
func GetCurrentUser(r *http.Request) *models.User {

	userValue := r.Context().Value("user")

	if userValue == nil {
		return nil
	}

	user, ok := userValue.(*models.User)
	if !ok {
		return nil
	}

	return user
}
