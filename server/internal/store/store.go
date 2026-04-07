package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sent_alerts (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			event_time TEXT    NOT NULL,
			alert_type TEXT    NOT NULL,
			sent_at    TEXT    NOT NULL,
			UNIQUE(event_time, alert_type)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) AlreadySent(eventTime time.Time, alertType string) (bool, error) {
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM sent_alerts WHERE event_time = ? AND alert_type = ?`,
		eventTime.UTC().Format(time.RFC3339), alertType,
	).Scan(&count)
	return count > 0, err
}

func (s *Store) MarkSent(eventTime time.Time, alertType string) error {
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO sent_alerts (event_time, alert_type, sent_at) VALUES (?, ?, ?)`,
		eventTime.UTC().Format(time.RFC3339),
		alertType,
		time.Now().UTC().Format(time.RFC3339),
	)
	return err
}
