CREATE TABLE account_balance (
    account_id  UUID PRIMARY KEY REFERENCES account(id),
    coin        INTEGER NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
