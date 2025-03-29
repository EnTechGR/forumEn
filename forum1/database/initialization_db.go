package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
)

// initDB initializes the database and returns a connection
func InitDB() (*sql.DB, error) {
	// Create database directory if it doesn't exist
	dbDir := "./database"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Connect to SQLite database (will be created if it doesn't exist)
	dbPath := filepath.Join(dbDir, "forum.db")
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Verify connection works
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Set some basic connection pool settings
	db.SetMaxOpenConns(10)

	// Initialize database schema and data
	if err := createTables(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	if err := createIndexes(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create indexes: %v", err)
	}

	if err := populateCategories(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to populate categories: %v", err)
	}

	return db, nil
}
