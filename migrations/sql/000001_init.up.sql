CREATE TABLE shopping_list
(
    user_id   uuid PRIMARY KEY NOT NULL UNIQUE,
    purchases JSONB            NOT NULL,
    version   integer          NOT NULL DEFAULT 1
);

CREATE TABLE inbox
(
    message_id uuid PRIMARY KEY NOT NULL UNIQUE
);
