--  create group_invited_users TABLE, migrate up

CREATE TABLE group_invited_users (
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    inviter_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (inviter_id) REFERENCES users (id)
    );