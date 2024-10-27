package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const operationName = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("#{%s}: open SQLite %w", operationName, err)
	}

	// TODO: Implement migrations

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("#{%s}: prepare statement: %w", operationName, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("#{%s}: exec statement: %w", operationName, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("#{%s}: prepare statement: %w", op, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("#{%s}: exec statement: %w", op, storage.ErrorAliasExists)
		}

		return fmt.Errorf("#{%s}: exec statement: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("#{%s}: prepare statement: %w", op, err)
	}

	var res string
	err = stmt.QueryRow(alias).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrorAliasNotFound
		}

		return "", fmt.Errorf("#{%s}: exec statement: %w", op, err)
	}

	return res, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("#{%s}: prepare statement: %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("#{%s}: exec statement: %w", op, err)
	}

	return nil
}
