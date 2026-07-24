package main

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func UpsertTLEs(db *sqlx.DB, tles []TLE) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	for _, sat := range tles {
		_, err := tx.NamedExec(
			`INSERT INTO TLEs
				(satname, line1, line2)
				VALUES
				(:satname, :line1, :line2)
				ON CONFLICT(satname)
				DO UPDATE SET
					line1 = excluded.line1,
					line2 = excluded.line2;`,
			&sat,
		)

		if err != nil {
			slog.Error("Failed to insert satellite", "error", err)
			return err
		}
	}

	slog.Info("TLEs inserted into DB")
	return tx.Commit()
}
