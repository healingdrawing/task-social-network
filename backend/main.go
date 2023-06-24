package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var (
	portHTTP string = "8000"
	fileDB   string = "./forum.db"
	db       *sql.DB
	reset    *bool
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Flags
	reset = flag.Bool("db-reset", false, "Reset database")
	flag.Parse()

	dbInit()
	defer db.Close()
	runMigrations(db)
	statementsCreation()

	// Static files, forum
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	registerHandlers()
	log.Println("starting forum at http://localhost:" + portHTTP + "/")
	log.Println("starting websocket at ws://localhost:" + portHTTP + "/ws")
	log.Fatal(http.ListenAndServe(":"+portHTTP, nil))
}
