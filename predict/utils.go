package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"strings"
)

type TLE struct {
	SatName string
	Line1   string
	Line2   string
}

func ParseTLEs(tleStr string) []TLE {
	reader := strings.NewReader(tleStr)
	scanner := bufio.NewScanner(reader)
	tles := []TLE{}
	tle := TLE{}

	for scanner.Scan() {
		tle.SatName = scanner.Text()
		scanner.Scan()
		tle.Line1 = scanner.Text()
		scanner.Scan()
		tle.Line2 = scanner.Text()
		tles = append(tles, tle)
	}

	if err := scanner.Err(); err != nil {
		slog.Error(fmt.Sprintf("Scanner error encountered parsing TLE string: %v\n", err))
	}

	return tles
}
