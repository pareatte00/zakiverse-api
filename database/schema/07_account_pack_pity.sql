CREATE TABLE account_pack_pity (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID        NOT NULL REFERENCES account(id),
    pack_id    UUID        NOT NULL REFERENCES pack(id) ON DELETE CASCADE,
    rarity     card_rarity NOT NULL,
    counter    INT         NOT NULL DEFAULT 0,
    UNIQUE(account_id, pack_id, rarity)
);

CREATE INDEX idx_account_pack_pity_account_pack ON account_pack_pity(account_id, pack_id);
