CREATE TYPE shopping_list_type as ENUM ('personal', 'shared');

CREATE TABLE shopping_lists
(
    shopping_list_id uuid PRIMARY KEY   NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    name             varchar(64)                        DEFAULT NULL,
    type             shopping_list_type NOT NULL        DEFAULT 'personal',
    owner_id         uuid               NOT NULL,
    purchases        jsonb              NOT NULL,
    recipe_names     jsonb              NOT NULL,
    version          integer            NOT NULL        DEFAULT 1
);

CREATE TABLE shopping_lists_users
(
    shopping_list_id uuid REFERENCES shopping_lists (shopping_list_id) NOT NULL UNIQUE,
    user_id          uuid                                              NOT NULL,
    accepted         boolean                                           NOT NULL DEFAULT false
);

CREATE TABLE inbox
(
    message_id uuid PRIMARY KEY         NOT NULL UNIQUE,
    timestamp  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()::timestamp
);
