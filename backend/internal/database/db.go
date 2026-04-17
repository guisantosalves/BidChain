package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("database.New: open: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database.New: ping: %w", err)
	}

	if err := runMigration(db); err != nil {
		return nil, fmt.Errorf("database.New: migration: %w", err)
	}

	return db, nil
}

func runMigration(db *sql.DB) error {
	migration, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("runMigration: Read file: %w", err)
	}

	_, err = db.Exec(string(migration))
	if err != nil {
		return fmt.Errorf("runMigration: Exec migration: %w", err)
	}

	return nil
}
