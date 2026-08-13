ALTER TABLE users
    ADD COLUMN IF NOT EXISTS plan_expires_at TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS payments (
    id                   BIGSERIAL PRIMARY KEY,
    user_id              BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    payload              TEXT NOT NULL UNIQUE,
    plan                 TEXT NOT NULL CHECK (plan IN ('plus', 'pro')),
    days                 INTEGER NOT NULL CHECK (days > 0),
    amount_kop           INTEGER NOT NULL CHECK (amount_kop > 0),
    currency             TEXT NOT NULL DEFAULT 'BYN',
    status               TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid')),
    telegram_charge_id   TEXT,
    paid_at              TIMESTAMPTZ,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS payments_user_id_idx ON payments (user_id);
