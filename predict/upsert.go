package main

import (
	"gopredict/api"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func UpsertTLEs(db *sqlx.DB, tles []api.TLE) error {
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

func UpsertStations(db *sqlx.DB, stns []api.Station) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	for _, stn := range stns {
		_, err := tx.NamedExec(
			`INSERT INTO Stations
				(stnname, latitude, longitude, altitude, minhorizon)
				VALUES
				(:stnname, :latitude, :longitude, :altitude, :minhorizon)
				ON CONFLICT(stnname)
				DO UPDATE SET
					latitude = excluded.latitude,
					longitude = excluded.longitude,
					altitude = excluded.altitude,
					minhorizon = excluded.minhorizon;`,
			&stn,
		)

		if err != nil {
			slog.Error("Failed to insert station", "error", err)
			return err
		}
	}

	slog.Info("Stations inserted into DB")
	return tx.Commit()
}
