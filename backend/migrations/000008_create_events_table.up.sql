--  create events TABLE, migrate up

CREATE TABLE events (
    id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    title VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    date DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (group_id) REFERENCES groups (id)
    );
