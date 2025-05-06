package state

import (
	"database/sql"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

// The state of the application.
// Currently just for checking our last attempt to read the RSS feed, and the last RSS feed published.
type State struct {
	// The underlying SQL connection
	db *sql.DB
}

func New(path string) (*State, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	setup := `create table if not exists "milkweed" (guid text primary key, published text);`
	_, err = db.Exec(setup)
	if err != nil {
		return nil, err
	}

	return &State{
		db: db,
	}, nil
}

func (s *State) Publish(guid string, t time.Time) error {
	iso8601 := t.Format(time.RFC3339)
	_, err := s.db.Exec(`insert or replace into milkweed (guid, published) values (?, ?)`, guid, iso8601)
	if err != nil {
		return err
	}
	return nil
}

func (s *State) IsPublished(guid string) (bool, error) {
	row := s.db.QueryRow("select count(*) from milkweed where guid = ?", guid)
	err := row.Err()
	if err != nil {
		return false, err
	}
	if row == nil {
		return false, nil
	}
	return true, nil
}
