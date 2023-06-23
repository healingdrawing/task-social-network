package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
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

	// DB
	dbInit()
	defer db.Close()

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

	statementsCreation()

	// Static files, forum
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	log.Println("starting forum at http://localhost:" + portHTTP + "/")
	// Websocket
	http.HandleFunc("/ws", wsConnection)
	log.Println("starting websocket at ws://localhost:" + portHTTP + "/ws")
	// API
	http.HandleFunc("/api/comment/submit", commentNewHandler)
	http.HandleFunc("/api/comment/get", commentGetHandler)

	http.HandleFunc("/api/chat/getmessages", chatMessagesHandler)
	http.HandleFunc("/api/chat/getusers", chatUsersHandler)
	http.HandleFunc("/api/chat/newmessage", chatNewHandler)
	http.HandleFunc("/api/chat/typing", chatTypingHandler)

	http.HandleFunc("/api/followrequest/reject", rejectFollowerHandler)
	http.HandleFunc("/api/followrequest/accept", approveFollowerHandler)
	http.HandleFunc("/api/followrequestlist", followRequestListHandler)

	http.HandleFunc("/api/post/get", postGetHandler)
	http.HandleFunc("/api/post/submit", postNewHandler)

	// todo: remove later. New HandleFunc for group posts/comments, and plural form looks weird for me f.e. api/group/post/get, because it looks like returns the list of posts. And above too for api/post/get @sagarishere
	http.HandleFunc("/api/group/submit", groupNewHandler)
	http.HandleFunc("/api/group/post/submit", groupPostNewHandler)
	http.HandleFunc("/api/group/posts/get", groupPostsGetHandler)
	http.HandleFunc("/api/group/comment/submit", groupCommentNewHandler)
	http.HandleFunc("/api/group/comments/get", groupCommentsGetHandler)
	http.HandleFunc("/api/group/join", groupJoinHandler)
	http.HandleFunc("/api/group/leave", groupLeaveHandler) // TODO: not part of audit, so untested
	http.HandleFunc("/api/group/invite", groupInviteHandler)
	http.HandleFunc("/api/group/invited", groupInvitedHandler)
	http.HandleFunc("/api/group/invite/accept", groupInviteAcceptHandler)
	http.HandleFunc("/api/group/invite/reject", groupInviteRejectHandler)
	http.HandleFunc("/api/group/requests", groupRequestsHandler)
	http.HandleFunc("/api/group/request/accept", groupRequestAcceptHandler)
	http.HandleFunc("/api/group/request/reject", groupRequestRejectHandler)

	http.HandleFunc("/api/user/check", sessionCheckHandler)
	http.HandleFunc("/api/user/following", FollowingHandler)
	http.HandleFunc("/api/user/followers", FollowersHandler)
	http.HandleFunc("/api/user/follow", FollowHandler)
	http.HandleFunc("/api/user/login", userLoginHandler)
	http.HandleFunc("/api/user/logout", userLogoutHandler)
	http.HandleFunc("/api/user/posts", userPostsHandler)
	http.HandleFunc("/api/user/privacy", changePrivacyHandler)
	http.HandleFunc("/api/user/profile", userProfileHandler)
	http.HandleFunc("/api/user/register", userRegisterHandler)
	http.HandleFunc("/api/user/unfollow", UnfollowHandler)
	// Server
	log.Fatal(http.ListenAndServe(":"+portHTTP, nil))
}
