--  create group_invited_users TABLE, migrate up

CREATE TABLE group_invited_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (id),
    group_id INTEGER NOT NULL REFERENCES groups (id),
    created_at DATETIME NOT NULL
    );