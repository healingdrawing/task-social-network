package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(db *sql.DB) {
	// get the current path of the directory in the os
	// and append the migrations folder to it
	currentPath, _ := os.Getwd()
	migrationsPath := "file://" + currentPath + "/migrations/"

	// Create a new instance of the sqlite3 driver
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Create a new instance of the migrate package
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		fileDB, driver)
	if err != nil {
		log.Printf("Could not create migrate instance: %v\n", err)
	}

	if err == nil {
		// New migration initialized
		// Migrate the database to the latest version
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Printf("Could not migrate: %v\n", err)
		}
		if err == migrate.ErrNoChange {
			log.Println("No migration needed")
		} else {
			log.Println("Migration successful")
		}
		// if migration is successful, then print the latest version
		version, _, _ := m.Version()
		log.Printf("Latest migrated version: %d\n", version)
	}
}
