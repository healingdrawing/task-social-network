-- make the post TABLE, migrate up

CREATE TABLE post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users (id),
    title VARCHAR NOT NULL,
    categories VARCHAR,
    content VARCHAR NOT NULL,
    privacy VARCHAR NOT NULL,
    picture BLOB,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );