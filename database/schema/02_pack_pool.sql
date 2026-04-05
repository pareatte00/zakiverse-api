CREATE TABLE pack_pool (
    id              UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100)  NOT NULL,
    description     TEXT,
    active_count    INT           NOT NULL DEFAULT 1,
    rotation_day    INT           NOT NULL DEFAULT 1,
    last_rotated_at TIMESTAMPTZ,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
