-- make the group_post_membership TABLE, migrate up

CREATE TABLE group_post_membership (
    group_id INTEGER,
    group_post_id INTEGER,
    PRIMARY KEY (group_id, group_post_id),
    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (group_post_id) REFERENCES group_posts (id)
    );