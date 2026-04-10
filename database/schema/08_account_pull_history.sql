CREATE TABLE account_pull_history (
    id          UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id  UUID          NOT NULL REFERENCES account(id),
    pack_id     UUID          NOT NULL REFERENCES pack(id) ON DELETE CASCADE,
    card_id     UUID          NOT NULL REFERENCES card(id),
    rarity      card_rarity   NOT NULL,
    is_pity     BOOLEAN       NOT NULL DEFAULT FALSE,
    is_featured BOOLEAN       NOT NULL DEFAULT FALSE,
    is_new      BOOLEAN       NOT NULL DEFAULT FALSE,
    pulled_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pull_history_account ON account_pull_history(account_id);
CREATE INDEX idx_pull_history_account_pack ON account_pull_history(account_id, pack_id);
