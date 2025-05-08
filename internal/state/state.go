package state

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

// The state of the application.
// Currently just for checking our last attempt to read the RSS feed, and the last RSS feed published.
type State struct {
	// The underlying SQL connection
	db     *sql.DB
	logger *log.Logger
	ctx    context.Context
}

func New(path string, logger *log.Logger, ctx context.Context) (*State, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	logger.Println("opened sqlite state at '", path, "'")

	setup := `create table if not exists "milkweed" (guid text primary key, published text);`
	_, err = db.Exec(setup)
	if err != nil {
		return nil, err
	}

	return &State{
		db:     db,
		logger: logger,
		ctx:    ctx,
	}, nil
}

func (s *State) Close() error {
	return s.db.Close()
}

func (s *State) Publish(guid string, t time.Time) error {
	rfc3339 := t.Format(time.RFC3339)
	_, err := s.db.Exec(`insert or replace into milkweed (guid, published) values (?, ?)`, guid, rfc3339)
	if err != nil {
		return err
	}
	return nil
}

func (s *State) IsPublished(guid string) (bool, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	var count int
	err := s.db.QueryRowContext(ctx, "select count(*) from milkweed where guid = ?", guid).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
