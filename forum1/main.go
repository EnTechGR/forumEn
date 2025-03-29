package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"forum/database"
	"forum/routes"
	"forum/utils"
)

func main() {
	// Load environment variables from .env file
	err := utils.LoadEnv(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Fetch SERVER_URL from environment variable
	serverUrl := os.Getenv("SERVER_URL")
	if serverUrl == "" {
		log.Fatal("SERVER_URL environment variable is not set")
	}

	// Parse the SERVER_URL to get the host
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		log.Fatalf("Invalid SERVER_URL: %v", err)
	}

	// Ensure the host is correctly formatted (including port if necessary)
	host := parsedUrl.Host
	if host == "" {
		// If the host part is empty, default to localhost:8080
		fmt.Println("Revert to default")
		host = "localhost:8080"

	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Setup routes
	handler := routes.SetupRoutes(db)

	// Start the server and log any fatal errors
	fmt.Printf("Server is running on %s\n", host)
	log.Fatal(http.ListenAndServe(host, handler))
}
