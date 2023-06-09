CREATE TYPE shopping_list_type as ENUM ('personal', 'shared');

CREATE TABLE shopping_lists
(
    shopping_list_id uuid PRIMARY KEY   NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    type             shopping_list_type NOT NULL        DEFAULT 'personal',
    owner_id         uuid               NOT NULL,
    purchases        jsonb              NOT NULL        DEFAULT '[]'::jsonb,
    version          integer            NOT NULL        DEFAULT 1 CHECK (version > 0)
);

CREATE INDEX shopping_lists_owner_id_key ON shopping_lists (owner_id);

CREATE TABLE shopping_lists_users
(
    shopping_list_id uuid REFERENCES shopping_lists (shopping_list_id) ON DELETE CASCADE NOT NULL,
    user_id          uuid                                                                NOT NULL,
    name             varchar(64) DEFAULT NULL,
    UNIQUE (shopping_list_id, user_id)
);

CREATE INDEX shopping_lists_users_user_id_key ON shopping_lists_users (user_id);

CREATE TABLE keys
(
    shopping_list_id uuid REFERENCES shopping_lists (shopping_list_id) ON DELETE CASCADE NOT NULL UNIQUE,
    key              uuid                                                                NOT NULL DEFAULT gen_random_uuid(),
    expires_at       TIMESTAMP WITH TIME ZONE                                            NOT NULL
);

CREATE TABLE inbox
(
    message_id uuid PRIMARY KEY         NOT NULL UNIQUE,
    timestamp  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()::timestamp
);
