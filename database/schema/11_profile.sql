CREATE TABLE profile (
    account_id      UUID PRIMARY KEY REFERENCES account(id),
    display_name    VARCHAR(50) NOT NULL,
    bio             TEXT,
    showcase_cards  JSONB NOT NULL DEFAULT '[]',
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_profile_display_name ON profile(display_name);
