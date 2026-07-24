package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Migrate(db *sqlx.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS TLEs (
		satname			TEXT UNIQUE PRIMARY KEY,
      	line1 			TEXT UNIQUE,
      	line2 	 		TEXT UNIQUE
    );`

	db.MustExec(schema)
}
