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
    );

    CREATE TABLE IF NOT EXISTS Stations (
    	stnname         TEXT UNIQUE PRIMARY KEY,
    	latitude        FLOAT NOT NULL,
    	longitude       FLOAT NOT NULL,
    	altitude        FLOAT NOT NULL,
    	minhorizon      FLOAT NOT NULL
     );`

	db.MustExec(schema)
}
