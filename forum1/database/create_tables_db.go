package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func createTables(db *sql.DB) error {
	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Define all table creation SQL statements
	tableStatements := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS user (
    		user_id TEXT PRIMARY KEY,
    		username TEXT NOT NULL UNIQUE 
        		CHECK (length(username) >= 3 AND length(username) <= 15)
        		CHECK (username GLOB '[a-zA-Z0-9_]*')
        		CHECK (username NOT GLOB '*[^a-zA-Z0-9_]*'),
    		email TEXT NOT NULL UNIQUE,
    		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,

		// User authentication table
		`CREATE TABLE IF NOT EXISTS user_auth (
    		user_id TEXT PRIMARY KEY,
    		password_hash TEXT NOT NULL 
        		CHECK (length(password_hash) = 60),
    		FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE
		);`,

		// Sessions table
		`CREATE TABLE IF NOT EXISTS sessions (
			user_id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL UNIQUE,
			ip_address TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE
		);`,

		// Categories table
		`CREATE TABLE IF NOT EXISTS categories (
    		category_id TEXT PRIMARY KEY,
    		category_number INTEGER NOT NULL UNIQUE,
    		category_name TEXT NOT NULL UNIQUE
		);`,

		// Posts table
		`CREATE TABLE IF NOT EXISTS posts (
			post_id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			category_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
		);`,

		// Comments table
		`CREATE TABLE IF NOT EXISTS comments (
			comment_id TEXT PRIMARY KEY,
			post_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE
		);`,

		// Reactions table
		`CREATE TABLE IF NOT EXISTS reactions (
    		reaction_id TEXT PRIMARY KEY,
    		user_id TEXT NOT NULL,
    		reaction_type INTEGER NOT NULL, -- 1 for like, 2 for dislike
    		comment_id TEXT,
    		post_id TEXT,
    		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    		FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    		FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE,
    		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    		CHECK ((post_id IS NULL AND comment_id IS NOT NULL) OR (post_id IS NOT NULL AND comment_id IS NULL))
		);`,
	}

	// Execute each table creation statement
	for _, stmt := range tableStatements {
		_, err = tx.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %s: %v", stmt, err)
		}
	}

	// Commit transaction
	return tx.Commit()
}
