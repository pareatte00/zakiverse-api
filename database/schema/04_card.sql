CREATE TYPE card_rarity AS ENUM ('common', 'uncommon', 'rare', 'epic', 'legendary');

CREATE TABLE card (
    id         UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    mal_id     INT           NOT NULL UNIQUE,
    anime_id   UUID          NOT NULL REFERENCES anime(id),
    rarity     card_rarity   NOT NULL,
    name       VARCHAR(100)  NOT NULL,
    image      VARCHAR(500)  NOT NULL,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_card_anime_id ON card(anime_id);
CREATE INDEX idx_card_rarity ON card(rarity);
