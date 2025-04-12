package state

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	_ "github.com/glebarez/go-sqlite"
)

// The state of the application.
// Currently just for checking our last attempt to read the RSS feed, and the last RSS feed published.
type State struct {
	// The latest published poast's date and time of publication
	Published *time.Time

	// The underlying SQL connection
	db *sql.DB
}

func New(path string) (*State, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	setup := `create table if not exists "milkweed" (published text);`
	_, err = db.Exec(setup)
	if err != nil {
		return nil, err
	}

	row := db.QueryRow("select count(*) from milkweed")
	var rowCount string
	err = row.Scan(&rowCount)
	if err != nil {
		return nil, err
	}
	count, err := strconv.Atoi(rowCount)
	if err != nil {
		return nil, err
	}
	var published time.Time
	if count == 0 {
		return &State{
			Published: &published,
			db:        db,
		}, nil
	}

	// get SQLite version
	row = db.QueryRow("select published from milkweed limit 1")
	var (
		rowPublished string
	)
	err = row.Scan(&rowPublished)
	if err != nil {
		return nil, err
	}
	published, err = dateparse.ParseAny(rowPublished)
	if err != nil {
		return nil, err
	}
	return &State{
		Published: &published,
		db:        db,
	}, nil
}

func (s *State) Publish(t time.Time) error {
	iso8601 := t.Format(time.RFC3339)
	_, err := s.db.Exec(`INSERT OR REPLACE INTO milkweed (published) VALUES (?)`, iso8601)
	if err != nil {
		return err
	}
	return nil
}
