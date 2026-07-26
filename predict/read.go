package main

import (
	"gopredict/api"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func ReadTLEBySatName(db *sqlx.DB, satname string) (*api.TLE, error) {
	rows, err := db.Queryx(
		`SELECT
			satname,
			line1,
			line2
        FROM TLEs
        WHERE satname=?
        LIMIT 1`,
		satname,
	)

	if err != nil {
		slog.Error("Failed to query TLEs table: ", "error", err)
		return nil, err
	}

	defer rows.Close()
	tle := api.TLE{}

	for rows.Next() {
		err = rows.StructScan(&tle)
		if err != nil {
			slog.Error("Failed to scan tle struct: ", "error", err)
		}
	}

	return &tle, nil
}

func ReadTLEs(db *sqlx.DB) ([]api.TLE, error) {
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
	tles := []api.TLE{}
	tle := api.TLE{}

	for rows.Next() {
		err = rows.StructScan(&tle)
		tles = append(tles, tle)
	}

	return tles, nil
}

func ReadStations(db *sqlx.DB) ([]api.Station, error) {
	rows, err := db.Queryx(
		`SELECT
			stnname,
			latitude,
			longitude,
			altitude,
			minhorizon
        FROM Stations`,
	)

	if err != nil {
		slog.Error("Failed to query Stations table: ", "error", err)
		return nil, err
	}

	defer rows.Close()
	stn := api.Station{}
	stns := []api.Station{}

	for rows.Next() {
		err = rows.StructScan(&stn)
		if err != nil {
			slog.Error("Failed to scan station struct: ", "error", err)
		}

		stns = append(stns, stn)
	}

	return stns, nil
}
