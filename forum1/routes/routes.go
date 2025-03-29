package routes

import (
	"database/sql"
	"net/http"

	"forum/handlers"
	"forum/middleware"
	"forum/repository"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(db *sql.DB) http.Handler {
	// Create repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Create services
	authService := handlers.NewAuthService(userRepo, sessionRepo)

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(sessionRepo, userRepo)

	// Create router (using standard net/http for simplicity)
	mux := http.NewServeMux()

	// Define web routes
	mux.HandleFunc("/", handlers.HomeHandler)

	// Define auth routes - public
	mux.HandleFunc("/api/auth/register", handlers.RegisterUser(authService))
	mux.HandleFunc("/api/auth/login", handlers.LoginUser(authService))

	// Protected routes - require authentication
	logoutHandler := authMiddleware.RequireAuth(http.HandlerFunc(handlers.LogoutUser(authService)))
	mux.Handle("/api/auth/logout", logoutHandler)

	// Apply the Authenticate middleware to all routes
	return authMiddleware.Authenticate(mux)
}
