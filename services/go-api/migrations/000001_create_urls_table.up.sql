CREATE TABLE urls (
    id           BIGSERIAL    PRIMARY KEY,
    code         TEXT         NOT NULL UNIQUE,
    original_url TEXT         NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    owner_id     BIGINT
);

-- Snabb lookup på kod — det här är den heta vägen vid redirect
CREATE INDEX idx_urls_code ON urls (code);
