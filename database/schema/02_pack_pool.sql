CREATE TYPE banner_type AS ENUM ('standard', 'featured', 'event', 'beginner', 'seasonal');
CREATE TYPE rotation_type AS ENUM ('none', 'weekly', 'monthly');
CREATE TYPE rotation_order_mode AS ENUM ('auto', 'manual');

CREATE TABLE pack_pool (
    id                  UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(100)        NOT NULL,
    description         TEXT,
    image               VARCHAR(500),
    banner_type         banner_type         NOT NULL DEFAULT 'standard',
    sort_order          INT                 NOT NULL DEFAULT 0,
    is_active           BOOLEAN             NOT NULL DEFAULT FALSE,
    open_at             TIMESTAMPTZ,
    close_at            TIMESTAMPTZ,
    active_count        INT                 NOT NULL DEFAULT 1,
    rotation_type       rotation_type       NOT NULL DEFAULT 'none',
    rotation_day        INT,
    rotation_interval   INT                 NOT NULL DEFAULT 1,
    rotation_hour       INT                 NOT NULL DEFAULT 0,
    rotation_order_mode rotation_order_mode NOT NULL DEFAULT 'auto',
    next_rotation_at    TIMESTAMPTZ,
    last_rotated_at     TIMESTAMPTZ,
    preview_days        INT                 NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ         NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pack_pool_is_active ON pack_pool(is_active);
CREATE INDEX idx_pack_pool_open_at ON pack_pool(open_at);
CREATE INDEX idx_pack_pool_banner_type ON pack_pool(banner_type);
CREATE INDEX idx_pack_pool_sort_order ON pack_pool(sort_order);
CREATE INDEX idx_pack_pool_next_rotation ON pack_pool(next_rotation_at);
