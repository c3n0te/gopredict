package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func UpdateSats() (string, error) {
	resp, err := http.Get("https://celestrak.org/NORAD/elements/gp.php?GROUP=cubesat&FORMAT=tle")
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to fetch URL: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error(fmt.Sprintf("Server returned bad status: %s", resp.Status))
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to read response body: %v", err))
		return "", err
	}

	tleStr := string(bodyBytes)
	return tleStr, nil
}
