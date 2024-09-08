package sqlite

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (storage *Storage, err error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLES IF NOT EXISTS url (
		id INTEGER PRIMARY KEY
		alias TEXT NOT NULL UNIQUE
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}

	return &Storage{db: db}, nil
}
