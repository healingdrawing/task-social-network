-- make the post TABLE, migrate up

CREATE TABLE post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users (id),
    title VARCHAR NOT NULL,
    categories VARCHAR,
    text VARCHAR NOT NULL,
    group_id INTEGER NOT NULL REFERENCES groups (id),
    privacy VARCHAR NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );