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
	portHTTP     string = "8080"
	fileDB       string = "./forum.db"
	db           *sql.DB
	reset        *bool
	CustomRouter *http.ServeMux
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

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, application/json, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			next.ServeHTTP(w, r)
		})
	}

	CustomRouter = http.NewServeMux()

	// Static files, forum
	CustomRouter.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	registerHandlers()
	log.Println("starting forum at http://localhost:" + portHTTP + "/")
	log.Println("starting websocket at ws://localhost:" + portHTTP + "/ws")
	log.Fatal(http.ListenAndServe(":"+portHTTP, corsMiddleware(CustomRouter)))
}
