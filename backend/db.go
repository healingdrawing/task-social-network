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
			text VARCHAR NOT NULL
			);
		CREATE TABLE message (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			from_id INTEGER NOT NULL REFERENCES users (id),
			to_id INTEGER NOT NULL REFERENCES users (id),
			text VARCHAR NOT NULL,
			time_sent DATETIME
			);
		CREATE TABLE groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			description TEXT,
			creator INTEGER,
			creation_date TIMESTAMP,
			privacy TEXT
			FOREIGN KEY (creator) REFERENCES members(id)
			);
		CREATE TABLE group_members (
			group_id INTEGER,
			member_id INTEGER,
			PRIMARY KEY (group_id, member_id),
			FOREIGN KEY (group_id) REFERENCES groups(id),
			FOREIGN KEY (member_id) REFERENCES members(id)
			);
		CREATE TABLE group_pending_members (
			group_id INTEGER,
			member_id INTEGER,
			PRIMARY KEY (group_id, member_id),
			FOREIGN KEY (group_id) REFERENCES groups(id),
			FOREIGN KEY (member_id) REFERENCES members(id)
			);
		CREATE TABLE group_images (
			group_id INTEGER,
			image_path TEXT,
			header_image_path TEXT,
			PRIMARY KEY (group_id, image),
			FOREIGN KEY (group_id) REFERENCES groups(id)
			);`)

		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println("DB reset")
	}
}

func statementsCreation() {
	for key, query := range map[string]string{
		"addUser":            `INSERT INTO users (username, age, gender, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?, ?, ?);`,
		"getAllUsers":        `SELECT id, username from users`,
		"getUserProfile":     `SELECT username, age, gender, first_name, last_name, email from users WHERE username=?`,
		"getUserbyID":        `SELECT username FROM users WHERE id = ?;`,
		"getUserID":          `SELECT id FROM users WHERE username = ?;`,
		"getUserCredentials": `SELECT username, password FROM users WHERE username = ? OR email = ?;`,

		"addSession":    `INSERT INTO session (uuid, user_id) VALUES (?, ?);`,
		"getSession":    `SELECT * FROM session WHERE uuid = ?;`,
		"getIDbyUUID":   `SELECT id FROM session INNER JOIN users ON users.id=user_id WHERE uuid = ?;`,
		"removeSession": `DELETE FROM session WHERE uuid = ?;`,

		"addPost":     `INSERT INTO post (user_id, title, categories, text) VALUES (?, ?, ?, ?);`,
		"getPosts":    `SELECT post.id, username, title, categories, text FROM post INNER JOIN users ON user_id=users.id ORDER BY post.id DESC;`,
		"addComment":  `INSERT INTO comment (user_id, post_id, text) VALUES (?, ?, ?);`,
		"getComments": `SELECT username, text FROM comment INNER JOIN users ON user_id = users.id WHERE post_id = ? ORDER BY comment.id DESC;`,
		"addMessage":  `INSERT INTO message (from_id, to_id, text, time_sent) VALUES (?, ?, ?, ?);`,
		"getMessages": `SELECT from_id, to_id, text, time_sent FROM message WHERE (from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?) ORDER BY time_sent DESC;`,

		"addGroup":  `INSERT INTO groups (name, description, creator, creation_date, privacy) VALUES (?, ?, ?, ?, ?);`,
		"getGroups": `SELECT id, name, description, creator, creation_date, privacy FROM groups ORDER BY creation_date DESC;`,
		"getGroup":  `SELECT id, name, description, creator, creation_date, privacy FROM groups WHERE id = ?;`,

		"addGroupMember":      `INSERT INTO group_members (group_id, member_id) VALUES (?, ?);`,
		"getGroupMembers":     `SELECT member_id FROM group_members WHERE group_id = ?;`,
		"getGroupMembersInfo": `SELECT username, first_name, last_name FROM users WHERE id = ?;`,

		"getGroupPendingMembers":     `SELECT member_id FROM group_pending_members WHERE group_id = ?;`,
		"addGroupPendingMember":      `INSERT INTO group_pending_members (group_id, member_id) VALUES (?, ?);`,
		"removeGroupPendingMember":   `DELETE FROM group_pending_members WHERE group_id = ? AND member_id = ?;`,
		"removeGroupMember":          `DELETE FROM group_members WHERE group_id = ? AND member_id = ?;`,
		"getGroupPendingMembersInfo": `SELECT username, first_name, last_name FROM users WHERE id = ?;`,

		"getGroupImages":         `SELECT image_path, header_image_path FROM group_images WHERE group_id = ?;`,
		"addGroupImage":          `INSERT INTO group_images (group_id, image_path, header_image_path) VALUES (?, ?, ?);`,
		"removeGroupImage":       `DELETE FROM group_images WHERE group_id = ? AND image_path = ?;`,
		"removeGroupHeaderImage": `DELETE FROM group_images WHERE group_id = ? AND header_image_path = ?;`,
		"getGroupPosts":          `SELECT post.id, username, title, categories, text FROM post INNER JOIN users ON user_id=users.id WHERE post.id IN (SELECT post_id FROM post_category WHERE category_id = ?) ORDER BY post.id DESC;`,
		"updateGroupImage":       `UPDATE group_images SET image_path = ? WHERE group_id = ? AND image_path = ?;`,
		"updateGroupHeaderImage": `UPDATE group_images SET header_image_path = ? WHERE group_id = ? AND header_image_path = ?;`,
	} {
		err := error(nil)
		statements[key], err = db.Prepare(query)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
