-- make the post TABLE, migrate up

CREATE TABLE post (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (id),
    title VARCHAR NOT NULL,
    categories VARCHAR,
    text VARCHAR NOT NULL
    );