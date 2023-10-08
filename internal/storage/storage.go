package storage

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %w", err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS passwords(
		id INTEGER PRIMARY KEY,
		url TEXT UNIQUE NOT NULL,
		alias TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_alias ON password(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare quuery statment: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("cannot execute query statment: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SavePassword(url, alias, password string) error {
	stmt, err := s.db.Prepare("INSERT INTO passwords(url, alias, password) VALUES(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("cannot prepare query: %w", err)
	}

	_, err = stmt.Exec(url, alias, password)
	if err != nil {
		return fmt.Errorf("cannot execute query: %w", err)
	}

	return nil
}

func (s *Storage) GetPassword(key string) (string, error) {
	if isUrl(key) {
		url := key

		rows, err := s.db.Query("SELECT password FROM passwords WHERE url = ?", url)
		if err != nil {
			return "", fmt.Errorf("cannot prepare query: %w", err)
		}
		defer rows.Close()

		var password string
		for rows.Next() {
			err := rows.Scan(&password)
			if err != nil {
				return "", fmt.Errorf("cannot get password string: %w", err)
			}
		}

		return password, nil
	} else {
		alias := key

		rows, err := s.db.Query("SELECT password FROM passwords WHERE alias = ?", alias)
		if err != nil {
			return "", fmt.Errorf("cannot prepare query: %w", err)
		}
		defer rows.Close()

		var password string
		for rows.Next() {
			err := rows.Scan(&password)
			if err != nil {
				return "", fmt.Errorf("cannot get password string: %w", err)
			}
		}

		return password, nil

	}
}

func (s *Storage) DeletePassword(key string) error {

	if isUrl(key) {
		url := key

		stmt, err := s.db.Prepare("DELETE FROM passwords WHERE url = ?")
		if err != nil {
			return fmt.Errorf("cannot prepare query: %w", err)
		}

		_, err = stmt.Exec(url)
		if err != nil {
			return fmt.Errorf("cannot execute query statement: %w", err)
		}

	} else {
		alias := key

		stmt, err := s.db.Prepare("DELETE FROM passwords WHERE alias = ?")
		if err != nil {
			return fmt.Errorf("cannot prepare query: %w", err)
		}

		_, err = stmt.Exec(alias)
		if err != nil {
			return fmt.Errorf("cannot execute query statement: %w", err)
		}
	}

	return nil
}

func (s *Storage) ResetPassword(key, newPassword string) error {
	if isUrl(key) {
		url := key

		stmt, err := s.db.Prepare("UPDATE passwords SET password = ? WHERE url = ?")
		if err != nil {
			return fmt.Errorf("cannot prepare query: %w", err)
		}

		_, err = stmt.Exec(newPassword, url)
		if err != nil {
			return fmt.Errorf("cannot execute query statement: %w", err)
		}
	} else {
		alias := key

		stmt, err := s.db.Prepare("UPDATE passwords SET password = ? WHERE alias = ?")
		if err != nil {
			return fmt.Errorf("cannot prepare query: %w", err)
		}

		_, err = stmt.Exec(newPassword, alias)
		if err != nil {
			return fmt.Errorf("cannot execute query statement: %w", err)
		}
	}

	return nil
}

func (s *Storage) GetAllPasswords() ([]string, error) {
	rows, err := s.db.Query("SELECT * FROM passwords")
	if err != nil {
		return nil, fmt.Errorf("failed to take query: %w", err)
	}

	resultRows := make([]string, 0)
	var id, url, alias, password string

	for rows.Next() {
		err := rows.Scan(&id, &url, &alias, &password)
		if err != nil {
			return nil, fmt.Errorf("failed to get rows: %w", err)
		}

		row := strings.Join([]string{id, url, alias, password}, " ")
		resultRows = append(resultRows, row)
	}

	return resultRows, nil
}
