-- create table session to migrate up
CREATE TABLE session (
    uuid VARCHAR PRIMARY KEY UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (id)
    );