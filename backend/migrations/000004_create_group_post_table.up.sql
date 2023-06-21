-- make the group_post TABLE, migrate up

CREATE TABLE group_post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users (id),
    title VARCHAR NOT NULL,
    categories VARCHAR,
    content VARCHAR NOT NULL,
    privacy VARCHAR NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );