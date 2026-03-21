CREATE TABLE account_card (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id  UUID        NOT NULL REFERENCES account(id),
    card_id     UUID        NOT NULL REFERENCES card(id),
    obtained_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (account_id, card_id)
);

CREATE INDEX idx_account_card_account_id ON account_card(account_id);
CREATE INDEX idx_account_card_card_id ON account_card(card_id);
