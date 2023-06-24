--  create events TABLE, migrate up

CREATE TABLE events (
    event_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    event_name VARCHAR NOT NULL,
    event_description VARCHAR NOT NULL,
    event_date DATETIME NOT NULL,
    event_location VARCHAR NOT NULL,
    creator_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY (event_id),
    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (creator_id) REFERENCES users (id)
    );
