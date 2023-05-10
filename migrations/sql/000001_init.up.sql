CREATE TABLE shopping_list
(
    user_id       uuid PRIMARY KEY NOT NULL UNIQUE,
    shopping_list JSONB            NOT NULL
);

CREATE TABLE inbox
(
    event_id  uuid PRIMARY KEY NOT NULL UNIQUE,
    type      VARCHAR(64)      NOT NULL,
    body      JSONB            NOT NULL DEFAULT '{}'::jsonb,
    processed BOOLEAN          NOT NULL DEFAULT false
);
