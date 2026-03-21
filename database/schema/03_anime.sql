CREATE TABLE anime (
    id          UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    mal_id      INT           NOT NULL UNIQUE,
    title       VARCHAR(255)  NOT NULL,
    synopsis    TEXT,
    cover_image VARCHAR(500),
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
