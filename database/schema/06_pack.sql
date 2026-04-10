CREATE TABLE pack (
    id                     UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    code                   VARCHAR(20)   NOT NULL DEFAULT '' UNIQUE,
    name                   VARCHAR(100)  NOT NULL,
    description            TEXT,
    image                  VARCHAR(500)  NOT NULL,
    name_image             VARCHAR(500),
    cards_per_pull         INT           NOT NULL DEFAULT 5,
    sort_order             INT           NOT NULL DEFAULT 0,
    config                 JSONB         NOT NULL DEFAULT '{}',
    pool_id                UUID          NOT NULL REFERENCES pack_pool(id) ON DELETE RESTRICT,
    rotation_order         INT,
    last_pool_activated_at TIMESTAMPTZ,
    created_at             TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE TABLE pack_card (
    id            UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    pack_id       UUID    NOT NULL REFERENCES pack(id) ON DELETE CASCADE,
    card_id       UUID    NOT NULL REFERENCES card(id) ON DELETE CASCADE,
    weight        FLOAT   NOT NULL DEFAULT 1.0,
    is_featured   BOOLEAN NOT NULL DEFAULT FALSE,
    featured_rate FLOAT,
    UNIQUE(pack_id, card_id)
);

CREATE INDEX idx_pack_card_pack_id ON pack_card(pack_id);
CREATE INDEX idx_pack_card_card_id ON pack_card(card_id);
CREATE INDEX idx_pack_pool_id ON pack(pool_id);
CREATE INDEX idx_pack_sort_order ON pack(sort_order);
