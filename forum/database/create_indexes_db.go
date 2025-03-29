package database

import (
	"database/sql"
	"fmt"
)

func createIndexes(db *sql.DB) error {
	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Define all index creation SQL statements
	indexStatements := []string{
		`CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_posts_category_id ON posts(category_id);`,
		`CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);`,
		`CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_reactions_user_id ON reactions(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_reactions_post_id ON reactions(post_id);`,
		`CREATE INDEX IF NOT EXISTS idx_reactions_comment_id ON reactions(comment_id);`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_session_id ON sessions(session_id);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_reactions_user_post ON reactions(user_id, post_id) WHERE post_id IS NOT NULL;`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_reactions_user_comment ON reactions(user_id, comment_id) WHERE comment_id IS NOT NULL;`,
	}

	// Execute each index creation statement
	for _, stmt := range indexStatements {
		_, err = tx.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %s: %v", stmt, err)
		}
	}

	// Commit transaction
	return tx.Commit()
}
