CREATE TABLE IF NOT EXISTS users (
    id         BIGSERIAL PRIMARY KEY,
    chat_id    BIGINT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS listings (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    UNIQUE (user_id, url)
);

CREATE INDEX IF NOT EXISTS listings_user_id_idx ON listings (user_id);
