package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"forum/config"
	"forum/models"
	"forum/repository"
	"forum/utils"
)

// AuthService handles authentication-related requests
type AuthService struct {
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
}

// AuthService creates a new AuthService
func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthService {
	return &AuthService{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}
}

// RegisterUser handles user registration
func RegisterUser(AuthService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse request body
		var reg models.UserRegistration
		err := json.NewDecoder(r.Body).Decode(&reg)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// validate username
		err = utils.ValidateUsername(reg.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// validate email
		err = utils.ValidateEmail(reg.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// validate password
		err = utils.ValidatePassword(reg.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create user
		user, err := AuthService.UserRepo.Create(reg)
		if err != nil {
			switch err {
			case config.ErrEmailTaken:
				http.Error(w, "Email is already taken", http.StatusConflict)
			case config.ErrUsernameTaken:
				http.Error(w, "Username is already taken", http.StatusConflict)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

// LoginUser handles user login
func LoginUser(AuthService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse request body
		var login models.UserLogin
		err := json.NewDecoder(r.Body).Decode(&login)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate request - basic required fields check
		if login.Email == "" || login.Password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Normalize email (convert to lowercase)
		login.Email = strings.ToLower(login.Email)

		// Basic email format validation
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(login.Email) {
			// Use generic error for security (don't reveal if email format is invalid)
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		user, err := AuthService.UserRepo.Authenticate(login)
		if err != nil {
			switch err {
			case config.ErrInvalidCredentials, config.ErrUserNotFound:
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Create a new session
		session, err := AuthService.SessionRepo.Create(user.ID, r.RemoteAddr)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		// Set the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    session.SessionID,
			Path:     "/",
			Expires:  session.ExpiresAt,
			HttpOnly: true,
			Secure:   r.TLS != nil, // Set Secure flag if TLS is enabled
			SameSite: http.SameSiteStrictMode,
		})

		// Return JSON response only
		w.Header().Set("Content-Type", "application/json")
		response := models.LoginResponse{
			User:      *user,
			SessionID: session.SessionID,
		}
		json.NewEncoder(w).Encode(response)
	}
}

// LogoutUser handles user logout
func LogoutUser(AuthService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get the session cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			// If no cookie, nothing to do
			w.WriteHeader(http.StatusOK)
			return
		}

		// Delete the session
		err = AuthService.SessionRepo.Delete(cookie.Value)
		if err != nil {
			http.Error(w, "Failed to logout", http.StatusInternalServerError)
			return
		}

		// Clear the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})

		w.WriteHeader(http.StatusOK)
	}
}
