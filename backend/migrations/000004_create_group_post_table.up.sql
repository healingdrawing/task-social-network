-- make the group_post TABLE, migrate up

CREATE TABLE group_post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title VARCHAR NOT NULL,
    categories VARCHAR,
    content VARCHAR NOT NULL,
    picture BLOB,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
    );