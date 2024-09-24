package postgresql

import (
	"database/sql"
	"fmt"
	"restapi/URL-Shortener/internal/storage"

	_ "github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (storage *Storage, err error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: Unique url
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url (
    id INTEGER PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL
	);`)

	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
`)

	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) { //del int64
	const op = "storage.postgresql.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	fmt.Println(err)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: Change when replacing db
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgresql.GetURL"

	if alias == "" {
		return "", fmt.Errorf("%s: alias id empty", op)
	}

	stmt, err := s.db.Prepare("SELECT url.url FROM url WHERE alias = ?")

	if err != nil {
		return "", fmt.Errorf("%s: failed to get db: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Query(alias)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrNoExtended(sqlite3.ErrNotFound) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: failed to get url: %w", op, err)
	}
	defer res.Close()

	// TODO: change to stmt.QueryRow().Scan(&name)
	var url string

	for res.Next() {
		err = res.Scan(&url)
		if err != nil {
			return "", fmt.Errorf("%s: failed to read a line: %w", op, err)
		}
	}

	return url, nil
}

func (s *Storage) DelURLByID(id int) (int, error) {
	op := "storage.postgresql.DelURLByID"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE id = ?")

	if err != nil {
		return 0, fmt.Errorf("%s: failed to get db: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrNoExtended(sqlite3.ErrNotFound) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return 0, fmt.Errorf("%s: failed to del url: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DelURLByAlias(alias string) (int, error) {
	op := "storage.postgresql.DelURLByAlias"

	stmt, err := s.db.Prepare("SELECT id FROM url WHERE alias = ?")

	if err != nil {
		return 0, fmt.Errorf("%s: failed to get db: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Query(alias)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrNoExtended(sqlite3.ErrNotFound) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return 0, fmt.Errorf("%s: failed to get id: %w", op, err)
	}

	var id int

	if !res.Next() {
		return 0, fmt.Errorf("%s: undefined alias: %s", op, alias)
	}

	for res.Next() {
		err = res.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("%s: failed to read a line: %w", op, err)
		}
	}

	_, err = s.db.Exec("DELETE FROM url WHERE alias = ?", alias)

	if err != nil {
		return 0, fmt.Errorf("%s: failed to del url: %w", op, err)
	}
	return id, nil
}
