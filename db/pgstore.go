package db

import (
	"database/sql"
	"fmt"

	"github.com/gchaincl/dotsql"
	_ "github.com/lib/pq" // The Postgres driver
)

// A PGStore implements storage against PostGres
type PGStore struct {
	conn *sql.DB
}

// NewPGStore connects to a postgres database
func NewPGStore(u, p, d string) (*PGStore, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		u, p, d)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(5)
	return &PGStore{conn: db}, nil
}

// ExecuteSchema runs all the sql commands int he given file to
// initialize the database
func (pg *PGStore) ExecuteSchema(filename string) error {
	dot, err := dotsql.LoadFromFile(filename)
	if err != nil {
		return err
	}
	for query := range dot.QueryMap() {
		dot.Exec(pg.conn, query)
	}
	return nil
}
