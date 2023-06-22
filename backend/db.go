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
			content VARCHAR NOT NULL,
			picture BLOB,
			created_at DATETIME NOT NULL
			);
		CREATE TABLE message (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			from_id INTEGER NOT NULL REFERENCES users (id),
			to_id INTEGER NOT NULL REFERENCES users (id),
			content VARCHAR NOT NULL,
			created_at DATETIME NOT NULL
			);
		CREATE TABLE groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			description TEXT,
			creator INTEGER,
			created_at TIMESTAMP,
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
		CREATE TABLE followers (
			user_id INTEGER,
			follower_id INTEGER,
			PRIMARY KEY (user_id, follower_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (follower_id) REFERENCES users(id)
			);
		CREATE TABLE followers_pending (
			user_id INTEGER,
			follower_id INTEGER,
			PRIMARY KEY (user_id, follower_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (follower_id) REFERENCES users(id)
			);
		CREATE TABLE post_category (
			post_id INTEGER NOT NULL,
			category_id INTEGER,
			PRIMARY KEY (post_id, category_id),
			FOREIGN KEY (post_id) REFERENCES post(id),
			FOREIGN KEY (category_id) REFERENCES category(id)
			);
		CREATE TABLE almost_private (
			user_id INTEGER,
			post_id INTEGER,
			PRIMARY KEY (user_id, post_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES post(id)
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
		"getUserIDByEmail":   `SELECT id FROM users WHERE email = ?;`,
		"getUserPrivacy":     `SELECT privacy FROM users WHERE id = ?;`,
		"getUserCredentials": `SELECT email, password FROM users WHERE email = ?;`,
		"updateUserPrivacy":  `UPDATE users SET privacy = ? WHERE id = ?;`,
		"getEmailByID":       `SELECT email FROM users WHERE id = ?;`,

		"addSession":    `INSERT INTO session (uuid, user_id) VALUES (?, ?);`,
		"getSession":    `SELECT * FROM session WHERE uuid = ?;`,
		"getIDbyUUID":   `SELECT id FROM session INNER JOIN users ON users.id=user_id WHERE uuid = ?;`,
		"getIDbyEmail":  `SELECT id FROM users WHERE email = ?;`,
		"removeSession": `DELETE FROM session WHERE uuid = ?;`,

		"addAlmostPrivate": `INSERT INTO almost_private (user_id, post_id) VALUES (?, ?);`,

		"addPost":     `INSERT INTO post (user_id, title, categories, content, privacy, picture, created_at) VALUES (?, ?, ?, ?, ?, ?, ?);`,
		"getPosts":    `SELECT post.id, title, content, categories, picture, first_name, last_name, email, created_at FROM post INNER JOIN users ON user_id=? ORDER BY created_at DESC;`,
		"addComment":  `INSERT INTO comment (user_id, post_id, content, picture, created_at) VALUES (?, ?, ?, ?, ?);`,
		"getComments": `SELECT first_name, last_name, content, picture FROM comment INNER JOIN users ON user_id = users.id WHERE post_id = ? ORDER BY comment.id DESC;`,
		"addMessage":  `INSERT INTO message (from_id, to_id, content, created_at) VALUES (?, ?, ?, ?);`,
		"getMessages": `SELECT from_id, to_id, content, created_at FROM message WHERE (from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?) ORDER BY created_at DESC;`,

		"addGroup":  `INSERT INTO groups (name, description, creator, created_at, privacy) VALUES (?, ?, ?, ?, ?);`,
		"getGroups": `SELECT id, name, description, creator, created_at, privacy FROM groups ORDER BY created_at DESC;`,
		"getGroup":  `SELECT id, name, description, creator, created_at, privacy FROM groups WHERE id = ?;`,

		"addGroupMember":           `INSERT INTO group_members (group_id, member_id) VALUES (?, ?);`,
		"getGroupMembers":          `SELECT member_id FROM group_members WHERE group_id = ?;`,
		"getGroupMembersInfo":      `SELECT nickname, first_name, last_name FROM users WHERE id = ?;`,
		"getGroupPendingMembers":   `SELECT member_id FROM group_pending_members WHERE group_id = ?;`,
		"addGroupPendingMember":    `INSERT INTO group_pending_members (group_id, member_id) VALUES (?, ?);`,
		"removeGroupPendingMember": `DELETE FROM group_pending_members WHERE group_id = ? AND member_id = ?;`,
		"removeGroupMember":        `DELETE FROM group_members WHERE group_id = ? AND member_id = ?;`,

		"addGroupInvitedUser":  `INSERT INTO group_invited_users (user_id, group_id, inviter_id, created_at) VALUES (?, ?, ?, ?);`,
		"getGroupInvitedUsers": `SELECT user_id, group_id, inviter_id, created_at FROM group_invited_users WHERE group_id = ?;`,

		"getFollowers":          `SELECT follower_id FROM followers WHERE user_id = ?;`,
		"getFollowersPending":   `SELECT follower_id FROM followers_pending WHERE user_id = ?;`,
		"addFollower":           `INSERT INTO followers (user_id, follower_id) VALUES (?, ?);`,
		"addFollowerPending":    `INSERT INTO followers_pending (user_id, follower_id) VALUES (?, ?);`,
		"removeFollower":        `DELETE FROM followers WHERE user_id = ? AND follower_id = ?;`,
		"removeFollowerPending": `DELETE FROM followers_pending WHERE user_id = ? AND follower_id = ?;`,
		"getFollowing":          `SELECT user_id FROM followers WHERE follower_id = ?;`,
		"doesSecondFollowFirst": `SELECT * FROM followers WHERE user_id = ? AND follower_id = ? LIMIT 1;`,

		"addGroupPost":           `INSERT INTO group_post (user_id, title, categories, content) VALUES (?, ?, ?, ?);`,
		"addGroupPostMembership": `INSERT INTO group_post_membership (group_id, group_post_id) VALUES (?, ?);`,
		"getGroupPosts":          `SELECT group_post.id, title, content, categories, first_name, last_name, email, created_at FROM group_post JOIN group_post_membership ON group_post.id = group_post_membership.group_post_id JOIN users ON group_post.user_id = users.id ORDER BY created_at DESC;`,

		"addGroupComment":  `INSERT INTO group_comment (user_id, group_post_id, content, picture) VALUES (?, ?, ?, ?);`,
		"getGroupComments": `SELECT email, first_name, last_name, nickname, content, picture FROM group_comment INNER JOIN users ON users.id = user_id WHERE group_post_id = ? ORDER BY group_comment.id DESC;`,
	} {
		err := error(nil)
		statements[key], err = db.Prepare(query)
		if err != nil {
			log.Print("Error preparing query: " + key)
			log.Fatal(err.Error())
		}
	}
}

// it was last added by @sagarishere, just moved it here, perhaps it better and will be used
// "getGroupPosts": `SELECT post.id, nickname, first_name, last_name, title, categories, content, created_at
// 							FROM post INNER JOIN users ON user_id=users.id WHERE post.id IN
// 							(SELECT post_id FROM post_category WHERE category_id = ?)
// 							ORDER BY created_at DESC;`,
