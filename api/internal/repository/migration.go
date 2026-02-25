package repository

import "database/sql"

func Migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS agents (
			id TEXT PRIMARY KEY,
			hostname TEXT NOT NULL,
			ip_address TEXT NOT NULL,
			port INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'active',
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS configs (
			id TEXT PRIMARY KEY,
			version INTEGER NOT NULL UNIQUE,
			data TEXT NOT NULL,
			poll_interval_seconds INTEGER NOT NULL DEFAULT 30,
			created_at DATETIME NOT NULL
		)`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
