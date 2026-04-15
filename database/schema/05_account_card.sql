CREATE TABLE account_card (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id  UUID        NOT NULL REFERENCES account(id),
    card_id     UUID        NOT NULL REFERENCES card(id),
    level       INTEGER     NOT NULL DEFAULT 1,
    obtained_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (account_id, card_id),
    CONSTRAINT account_card_level_check CHECK (level >= 1 AND level <= 5)
);

CREATE INDEX idx_account_card_account_id ON account_card(account_id);
CREATE INDEX idx_account_card_card_id ON account_card(card_id);
