package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	file, err := os.Create("./logs/gopredict.log")
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	handler := slog.NewJSONHandler(file, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	db, err := sqlx.Connect("sqlite3", "./gopredict.db")
	if err != nil {
		slog.Error("Failed to open SQLite Database: ", "error", err)
		return
	}
	defer db.Close()

	Migrate(db)
	p := tea.NewProgram(InitialModel(db))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
