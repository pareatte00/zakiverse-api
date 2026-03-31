CREATE TABLE pack (
    id             UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(100)  NOT NULL,
    description    TEXT,
    image          VARCHAR(500)  NOT NULL,
    cards_per_pull INT           NOT NULL DEFAULT 5,
    is_active      BOOLEAN       NOT NULL DEFAULT FALSE,
    open_at        TIMESTAMPTZ,
    close_at       TIMESTAMPTZ,
    config         JSONB         NOT NULL DEFAULT '{}',
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE TABLE pack_card (
    id       UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    pack_id  UUID    NOT NULL REFERENCES pack(id) ON DELETE CASCADE,
    card_id  UUID    NOT NULL REFERENCES card(id) ON DELETE CASCADE,
    weight   FLOAT   NOT NULL DEFAULT 1.0,
    UNIQUE(pack_id, card_id)
);

CREATE INDEX idx_pack_card_pack_id ON pack_card(pack_id);
CREATE INDEX idx_pack_card_card_id ON pack_card(card_id);
CREATE INDEX idx_pack_is_active ON pack(is_active);
CREATE INDEX idx_pack_open_at ON pack(open_at);
