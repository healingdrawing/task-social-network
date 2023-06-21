-- make the group_comment TABLE, migrate up

CREATE TABLE group_comment (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (id),
    group_post_id INTEGER NOT NULL REFERENCES group_post (id),
    content VARCHAR NOT NULL,
    create_time DATETIME NOT NULL
    );