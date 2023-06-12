package tests

import (
	_ "embed"
	"github.com/jmoiron/sqlx"
	"github.com/pavel-one/sensors/internal/sql"
)

//go:embed mocks/sensors.sql
var InsertSensors string

func GetTestDB() *sqlx.DB {
	db, err := sql.Connect("test")
	if err != nil {
		return nil
	}

	_, err = db.Exec(InsertSensors)
	if err != nil {
		return nil
	}

	return db
}