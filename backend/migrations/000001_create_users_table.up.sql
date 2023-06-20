--create users table for migration up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    dob VARCHAR NOT NULL,
    avatar BLOB,
    nickname VARCHAR,
    about_me VARCHAR,
    privacy TEXT NOT NULL
    );