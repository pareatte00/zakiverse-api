CREATE TYPE check_in_type AS ENUM ('recurring', 'streak', 'calendar');
CREATE TYPE check_in_reset_policy AS ENUM ('rolling', 'daily_reset', 'weekly_reset', 'monthly_reset');

CREATE TABLE check_in_plan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(50) NOT NULL UNIQUE,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    type            check_in_type NOT NULL,
    interval        INTEGER NOT NULL,
    max_claims      INTEGER NOT NULL DEFAULT 0,
    rewards         JSONB NOT NULL,
    reset_policy    check_in_reset_policy NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    starts_at       TIMESTAMPTZ,
    ends_at         TIMESTAMPTZ,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE check_in_record (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id      UUID NOT NULL REFERENCES account(id),
    plan_id         UUID NOT NULL REFERENCES check_in_plan(id),
    claim_count     INTEGER NOT NULL DEFAULT 0,
    streak          INTEGER NOT NULL DEFAULT 0,
    last_claimed    TIMESTAMPTZ,
    reset_at        TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(account_id, plan_id)
);

CREATE INDEX idx_check_in_record_account ON check_in_record(account_id);
