CREATE TYPE account_role AS ENUM ('user', 'admin');

CREATE TABLE account (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    username   VARCHAR(50)  NOT NULL UNIQUE,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    role       account_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
