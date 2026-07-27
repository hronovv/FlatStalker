CREATE TABLE IF NOT EXISTS seen_ads (
    listing_id BIGINT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    ad_id      BIGINT NOT NULL,
    seen_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (listing_id, ad_id)
);

CREATE INDEX IF NOT EXISTS seen_ads_listing_id_idx ON seen_ads (listing_id);
