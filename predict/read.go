package main

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func ReadTLEs(db *sqlx.DB) ([]TLE, error) {
	rows, err := db.Queryx(
		`SELECT
			satname,
			line1,
			line2
        FROM TLEs`,
	)

	if err != nil {
		slog.Error("Failed to query TLEs table: ", "error", err)
		return nil, err
	}

	defer rows.Close()
	tles := []TLE{}
	tle := TLE{}

	for rows.Next() {
		err = rows.StructScan(&tle)
		tles = append(tles, tle)
	}

	return tles, nil
}
