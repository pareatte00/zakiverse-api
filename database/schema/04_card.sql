CREATE TABLE card_tag (
    id         UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(50)   NOT NULL,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_card_tag_name ON card_tag(LOWER(name));

CREATE TYPE card_rarity AS ENUM ('common', 'uncommon', 'rare', 'epic', 'legendary', 'prismatic');

CREATE TABLE card (
    id         UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    mal_id     INT           NOT NULL,
    anime_id   UUID          NOT NULL REFERENCES anime(id),
    rarity     card_rarity   NOT NULL,
    name       VARCHAR(100)  NOT NULL,
    image      VARCHAR(500)  NOT NULL,
    config     JSONB         NOT NULL DEFAULT '{}',
    tag_id     UUID          NOT NULL REFERENCES card_tag(id),
    favorite   INT           NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_card_anime_id ON card(anime_id);
CREATE INDEX idx_card_rarity ON card(rarity);
CREATE UNIQUE INDEX idx_card_mal_id_tag_id ON card(mal_id, tag_id);
