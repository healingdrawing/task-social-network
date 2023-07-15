--  create events TABLE, migrate up

CREATE TABLE event_participants (
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    decision VARCHAR NOT NULL,
    PRIMARY KEY (event_id, user_id),
    FOREIGN KEY (event_id) REFERENCES events (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
    );
