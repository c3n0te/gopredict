package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopredict/api"
	"log/slog"
	"os"
	"strings"
)

func ParseTLEs(tleStr string) []api.TLE {
	reader := strings.NewReader(tleStr)
	scanner := bufio.NewScanner(reader)
	tles := []api.TLE{}
	tle := api.TLE{}

	for scanner.Scan() {
		tle.SatName = scanner.Text()
		scanner.Scan()
		tle.Line1 = scanner.Text()
		scanner.Scan()
		tle.Line2 = scanner.Text()
		tles = append(tles, tle)
	}

	if err := scanner.Err(); err != nil {
		slog.Error(fmt.Sprintf("Scanner error encountered parsing api.TLE string: %v\n", err))
	}

	return tles
}

func ParseStationFile() ([]api.Station, error) {
	stnData, err := os.ReadFile("./data/stations.json")
	if err != nil {
		return nil, err
	}

	var stns []api.Station
	json.Unmarshal(stnData, &stns)
	return stns, nil
}
