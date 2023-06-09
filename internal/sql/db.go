package sql

import (
	_ "embed" // embed database schema
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // database driver
)

//go:embed db/schema.sql
var schema string

// Connect to database
func Connect(dbname string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", fmt.Sprintf("./%s.sqlite3", dbname))
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return db, nil
}
