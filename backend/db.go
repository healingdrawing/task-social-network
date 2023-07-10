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
			creator_id INTEGER,
			created_at TIMESTAMP,
			FOREIGN KEY (creator_id) REFERENCES users (id)
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

// todo: CHECK IT! getPostsAbleToSee cybermonster is not checked properly
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

		"addPost":  `INSERT INTO post (user_id, title, categories, content, privacy, picture, created_at) VALUES (?, ?, ?, ?, ?, ?, ?);`,
		"getPosts": `SELECT post.id, title, content, categories, picture, post.privacy, created_at, email, first_name, last_name FROM post INNER JOIN users ON post.user_id = users.id WHERE post.user_id = ? ORDER BY created_at DESC;`,

		"getPostsAbleToSee": `SELECT post.id, title, content, categories, picture, post.privacy, created_at, email, first_name, last_name FROM post INNER JOIN users ON users.id = post.user_id LEFT JOIN followers ON followers.user_id = post.user_id AND followers.follower_id = ? LEFT JOIN almost_private ON almost_private.user_id = ? AND almost_private.post_id = post.id WHERE post.user_id = ? OR post.privacy = "public" OR (post.privacy = "private" AND followers.follower_id IS NOT NULL) OR (almost_private.post_id IS NOT NULL) ORDER BY created_at DESC;`,

		"getPostsAbleToSeeToVisitor": `SELECT post.id, title, content, categories, picture, post.privacy, created_at, email, first_name, last_name FROM post INNER JOIN users ON users.id = post.user_id WHERE (post.user_id = ? AND ? = ? OR post.privacy = "public" AND post.user_id = ? OR (post.privacy = "private" AND post.user_id = ? AND EXISTS (SELECT 1 FROM followers WHERE followers.user_id = ? AND followers.follower_id = ?)) OR post.user_id = ? AND EXISTS (SELECT 1 FROM almost_private WHERE almost_private.post_id = post.id AND almost_private.user_id = ?)) ORDER BY created_at DESC;`,

		"addComment":  `INSERT INTO comment (user_id, post_id, content, picture, created_at) VALUES (?, ?, ?, ?, ?);`,
		"getComments": `SELECT email, first_name, last_name, content, picture, created_at FROM comment INNER JOIN users ON user_id = users.id WHERE post_id = ? ORDER BY comment.id DESC;`,
		"addMessage":  `INSERT INTO message (from_id, to_id, content, created_at) VALUES (?, ?, ?, ?);`,
		"getMessages": `SELECT from_id, to_id, content, created_at FROM message WHERE (from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?) ORDER BY created_at DESC;`,

		"addGroup":     `INSERT INTO groups (name, description, creator_id, created_at) VALUES (?, ?, ?, ?);`,
		"getGroup":     `SELECT id, name, description, creator_id, created_at FROM groups WHERE id = ?;`,
		"getAllGroups": `SELECT groups.id, name, description, created_at, email, first_name, last_name FROM groups INNER JOIN users ON users.id = creator_id ORDER BY created_at DESC;`,
		"getGroups":    `SELECT groups.id, name, description, created_at, email, first_name, last_name FROM groups INNER JOIN users ON users.id = creator_id INNER JOIN group_members ON group_id = groups.id WHERE member_id = ? ORDER BY created_at DESC;`,

		"getCreatorAllGroupsPendings": `SELECT group_pending_members.group_id, group_pending_members.member_id, groups.name, groups.description, users.email, users.first_name, users.last_name FROM groups INNER JOIN group_pending_members ON group_pending_members.group_id = groups.id INNER JOIN users ON group_pending_members.member_id = users.id WHERE groups.creator_id = ?`,

		"addGroupMember":  `INSERT INTO group_members (group_id, member_id) VALUES (?, ?);`,
		"getGroupMembers": `SELECT member_id FROM group_members WHERE group_id = ?;`,
		"getGroupMember":  `SELECT member_id FROM group_members WHERE group_id = ? AND member_id = ?;`,

		"getGroupMembersInfo":    `SELECT nickname, first_name, last_name FROM users WHERE id = ?;`,
		"getGroupPendingMembers": `SELECT member_id FROM group_pending_members WHERE group_id = ?;`,
		"getGroupPendingMember":  `SELECT member_id FROM group_pending_members WHERE group_id = ? AND member_id = ?;`,

		"addGroupPendingMember":    `INSERT INTO group_pending_members (group_id, member_id) VALUES (?, ?);`,
		"removeGroupPendingMember": `DELETE FROM group_pending_members WHERE group_id = ? AND member_id = ?;`,
		"removeGroupMember":        `DELETE FROM group_members WHERE group_id = ? AND member_id = ?;`,

		"addGroupInvitedUser": `INSERT INTO group_invited_users (user_id, group_id, inviter_id, created_at) VALUES (?, ?, ?, ?);`,

		"getUserInvites": `SELECT groups.id, groups.name, groups.description, group_invited_users.created_at, users.email, users.first_name, users.last_name FROM group_invited_users JOIN groups ON groups.id = group_invited_users.group_id JOIN users ON users.id = group_invited_users.inviter_id WHERE group_invited_users.user_id = ?;`,

		"getGroupInvitedUsers":   `SELECT user_id, group_id, inviter_id, created_at FROM group_invited_users WHERE group_id = ?;`,
		"removeGroupInvitedUser": `DELETE FROM group_invited_users WHERE user_id = ? AND group_id = ?;`,

		"addGroupPost":           `INSERT INTO group_post (user_id, title, categories, content, picture, created_at) VALUES (?, ?, ?, ?, ?, ?);`,
		"addGroupPostMembership": `INSERT INTO group_post_membership (group_id, group_post_id) VALUES (?, ?);`,
		"getGroupPosts":          `SELECT group_post.id, title, content, categories, first_name, last_name, email, created_at, picture FROM group_post JOIN group_post_membership ON group_post.id = group_post_membership.group_post_id JOIN users ON group_post.user_id = users.id ORDER BY created_at DESC;`,

		"addGroupComment":  `INSERT INTO group_comment (user_id, group_post_id, content, picture, created_at) VALUES (?, ?, ?, ?, ?);`,
		"getGroupComments": `SELECT email, first_name, last_name, content, picture, created_at FROM group_comment INNER JOIN users ON users.id = user_id WHERE group_post_id = ? ORDER BY group_comment.id DESC;`,

		"addEvent":                  `INSERT INTO events (group_id, event_name, event_description, event_date, event_location, creator_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?);`,
		"addEventParticipant":       `INSERT INTO event_participants (event_id, user_id, status, status_updated_at) VALUES (?, ?, ?, ?);`,
		"getEvents":                 `SELECT event_id, group_id, event_name, event_description, event_date, event_location, creator_id, created_at FROM events WHERE group_id = ? ORDER BY created_at DESC;`,
		"getEvent":                  `SELECT event_id, group_id, event_name, event_description, event_date, event_location, creator_id, created_at FROM events WHERE event_id = ? LIMIT 1;`,
		"getEventParticipants":      `SELECT user_id, first_name, last_name, email, avatar, status, status_updated_at FROM event_participants INNER JOIN users ON users.id = user_id WHERE event_id = ? ORDER BY status_updated_at DESC;`,
		"updateEventParticipant":    `UPDATE event_participants SET status = ?, status_updated_at = ? WHERE event_id = ? AND user_id = ?;`,
		"getEventParticipantStatus": `SELECT status FROM event_participants WHERE event_id = ? AND user_id = ? LIMIT 1;`,
		"getUserIDwithEventCount":   `SELECT COUNT(*) FROM event_participants WHERE event_id = ? AND user_id = ?;`,

		"getFollowers":                   `SELECT follower_id FROM followers WHERE user_id = ?;`,
		"getFollowersPending":            `SELECT follower_id FROM followers_pending WHERE user_id = ?;`,
		"addFollower":                    `INSERT INTO followers (user_id, follower_id) VALUES (?, ?);`,
		"addFollowerPending":             `INSERT INTO followers_pending (user_id, follower_id) VALUES (?, ?);`,
		"removeFollower":                 `DELETE FROM followers WHERE user_id = ? AND follower_id = ?;`,
		"removeFollowerPending":          `DELETE FROM followers_pending WHERE user_id = ? AND follower_id = ?;`,
		"getFollowing":                   `SELECT user_id FROM followers WHERE follower_id = ?;`,
		"doesSecondFollowFirst":          `SELECT * FROM followers WHERE user_id = ? AND follower_id = ? LIMIT 1;`,
		"doesSecondRequesterFollowFirst": `SELECT * FROM followers_pending WHERE user_id = ? AND follower_id = ? LIMIT 1;`,
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
