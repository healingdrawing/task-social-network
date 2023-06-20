package main

import (
	"database/sql"
	"log"
	"os"
)

var statements = map[string]*sql.Stmt{}

func dbInit() {

	_, err := os.Stat(fileDB)

	if os.IsNotExist(err) {
		*reset = true
	}

	db, err = sql.Open("sqlite3", fileDB)
	if err != nil {
		log.Fatal(err)
	}
	if *reset {
		_, err := db.Exec(`
		DROP TABLE IF EXISTS users;
		DROP TABLE IF EXISTS session;
		DROP TABLE IF EXISTS post;
		DROP TABLE IF EXISTS category;
		DROP TABLE IF EXISTS post_category;
		DROP TABLE IF EXISTS comment;
		DROP TABLE IF EXISTS message;
		
		CREATE TABLE comment (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			user_id INTEGER NOT NULL REFERENCES users (id),
			post_id INTEGER NOT NULL REFERENCES post (id),
			text VARCHAR NOT NULL,
			create_time DATETIME NOT NULL
			);
		CREATE TABLE message (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			from_id INTEGER NOT NULL REFERENCES users (id),
			to_id INTEGER NOT NULL REFERENCES users (id),
			text VARCHAR NOT NULL,
			time_sent DATETIME NOT NULL
			);
		CREATE TABLE groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			description TEXT,
			creator INTEGER,
			creation_date TIMESTAMP,
			privacy TEXT,
			FOREIGN KEY (creator) REFERENCES users (id)
			);
		CREATE TABLE group_members (
			group_id INTEGER,
			member_id INTEGER,
			PRIMARY KEY (group_id, member_id),
			FOREIGN KEY (group_id) REFERENCES groups(id),
			FOREIGN KEY (member_id) REFERENCES users(id)
			);
		CREATE TABLE group_pending_members (
			group_id INTEGER,
			member_id INTEGER,
			PRIMARY KEY (group_id, member_id),
			FOREIGN KEY (group_id) REFERENCES groups(id),
			FOREIGN KEY (member_id) REFERENCES users(id)
			);
		CREATE TABLE category (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
			);
		CREATE TABLE post_category (
			post_id INTEGER NOT NULL,
			category_id INTEGER,
			PRIMARY KEY (post_id, category_id),
			FOREIGN KEY (post_id) REFERENCES post(id),
			FOREIGN KEY (category_id) REFERENCES category(id)
			);`)

		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println("DB reset")
	}
}

func statementsCreation() {
	for key, query := range map[string]string{
		"addUser":            `INSERT INTO users (email, password, first_name, last_name, dob, avatar, nickname, about_me, privacy) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`,
		"getAllUsers":        `SELECT id, email, first_name, last_name, nickname from users;`,
		"getUserProfile":     `SELECT email, first_name, last_name, dob, avatar, nickname, about_me, privacy from users WHERE id = ?;`,
		"getUserbyID":        `SELECT email, first_name, last_name, nickname FROM users WHERE id = ?;`,
		"getUserID":          `SELECT id FROM users WHERE email = ?;`,
		"getUserCredentials": `SELECT email, password FROM users WHERE email = ?;`,

		"addSession":    `INSERT INTO session (uuid, user_id) VALUES (?, ?);`,
		"getSession":    `SELECT * FROM session WHERE uuid = ?;`,
		"getIDbyUUID":   `SELECT id FROM session INNER JOIN users ON users.id=user_id WHERE uuid = ?;`,
		"removeSession": `DELETE FROM session WHERE uuid = ?;`,

		"addPost":     `INSERT INTO post (user_id, title, categories, text) VALUES (?, ?, ?, ?);`,
		"getPosts":    `SELECT post.id, first_name, last_name, nickname, title, categories, text FROM post INNER JOIN users ON user_id=users.id ORDER BY post.id DESC;`,
		"addComment":  `INSERT INTO comment (user_id, post_id, text) VALUES (?, ?, ?);`,
		"getComments": `SELECT nickname, text FROM comment INNER JOIN users ON user_id = users.id WHERE post_id = ? ORDER BY comment.id DESC;`,
		"addMessage":  `INSERT INTO message (from_id, to_id, text, time_sent) VALUES (?, ?, ?, ?);`,
		"getMessages": `SELECT from_id, to_id, text, time_sent FROM message WHERE (from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?) ORDER BY time_sent DESC;`,

		"addGroup":  `INSERT INTO groups (name, description, creator, creation_date, privacy) VALUES (?, ?, ?, ?, ?);`,
		"getGroups": `SELECT id, name, description, creator, creation_date, privacy FROM groups ORDER BY creation_date DESC;`,
		"getGroup":  `SELECT id, name, description, creator, creation_date, privacy FROM groups WHERE id = ?;`,

		"addGroupMember":      `INSERT INTO group_members (group_id, member_id) VALUES (?, ?);`,
		"getGroupMembers":     `SELECT member_id FROM group_members WHERE group_id = ?;`,
		"getGroupMembersInfo": `SELECT nickname, first_name, last_name FROM users WHERE id = ?;`,

		"getGroupPendingMembers":     `SELECT member_id FROM group_pending_members WHERE group_id = ?;`,
		"addGroupPendingMember":      `INSERT INTO group_pending_members (group_id, member_id) VALUES (?, ?);`,
		"removeGroupPendingMember":   `DELETE FROM group_pending_members WHERE group_id = ? AND member_id = ?;`,
		"removeGroupMember":          `DELETE FROM group_members WHERE group_id = ? AND member_id = ?;`,
		"getGroupPendingMembersInfo": `SELECT nickname, first_name, last_name FROM users WHERE id = ?;`,

		"getGroupPosts": `SELECT post.id, nickname, first_name, last_name, title, categories, text, created_at 
							FROM post INNER JOIN users ON user_id=users.id WHERE post.id IN 
							(SELECT post_id FROM post_category WHERE category_id = ?) 
							ORDER BY created_at DESC;`,
	} {
		err := error(nil)
		statements[key], err = db.Prepare(query)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
