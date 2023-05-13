CREATE TABLE shopping_list
(
    user_id       uuid PRIMARY KEY NOT NULL UNIQUE,
    shopping_list JSONB            NOT NULL,
    version       integer          NOT NULL DEFAULT 1
);

CREATE TABLE inbox
(
    message_id uuid PRIMARY KEY NOT NULL UNIQUE
);
