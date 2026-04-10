CREATE TABLE card_tag (
    id         UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(50)   NOT NULL UNIQUE,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE TYPE card_rarity AS ENUM ('common', 'uncommon', 'rare', 'epic', 'legendary', 'prismatic');

CREATE TABLE card (
    id         UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    mal_id     INT           NOT NULL,
    anime_id   UUID          NOT NULL REFERENCES anime(id),
    rarity     card_rarity   NOT NULL,
    name       VARCHAR(100)  NOT NULL,
    image      VARCHAR(500)  NOT NULL,
    config     JSONB         NOT NULL DEFAULT '{}',
    tag_id     UUID          REFERENCES card_tag(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_card_anime_id ON card(anime_id);
CREATE INDEX idx_card_rarity ON card(rarity);
