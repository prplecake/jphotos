package main

import (
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

func dbMigrate(command, dbURL string, version int) {
	m, err := migrate.New(
		"file://./db/migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal("migrate.New():", err)
	}
	currentVersion, dirty, err := m.Version()
	if err != nil {
		if !strings.Contains(err.Error(), "no migration") {
			// Assume the database is freshly created
			log.Fatalf("error getting migration version information: %s", err.Error())
		}
	}
	switch command {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal("m.Up():", err)
		}
		log.Print("Migrate Up: Success")
	case "down":
		if err := m.Migrate(currentVersion - 1); err != nil {
			log.Fatal("m.Migrate():", err)
		}
		log.Print("Migrate Down: Success")
	case "force":
		if err := m.Force(version); err != nil {
			log.Fatal("m.Force():", err)
		}
		log.Print("Migrate Force: Success")
	case "status":
		log.Print("Current migration version: ", currentVersion)
		log.Print("Migration dirty? ", dirty)
	}
}
