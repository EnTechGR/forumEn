package database

import (
	"database/sql"
	"fmt"
	"forum/utils"
)

func populateCategories(db *sql.DB) error {
	// Check if categories already exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing categories: %v", err)
	}

	// If categories already exist, skip population
	if count > 0 {
		return nil
	}

	// Define the categories to add
	categories := []string{
		"General Discussion",
		"Programming",
		"Golang",
		"Web Development",
		"Database Systems",
		"DevOps",
		"Mobile Development",
		"Machine Learning",
		"Security",
		"Off-Topic",
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the insert statement - now includes category_id and category_number
	stmt, err := tx.Prepare("INSERT INTO categories (category_id, category_number, category_name) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Insert each category with an auto-incremented number starting from 1
	for i, category := range categories {
		categoryID := utils.GenerateUUID()
		categoryNumber := i + 1

		_, err = stmt.Exec(categoryID, categoryNumber, category)
		if err != nil {
			return fmt.Errorf("failed to insert category '%s': %v", category, err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
