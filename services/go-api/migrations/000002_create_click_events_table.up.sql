CREATE TABLE click_events (
    id           BIGSERIAL    PRIMARY KEY,
    code         TEXT         NOT NULL,
    original_url TEXT         NOT NULL,
    clicked_at   TIMESTAMPTZ  NOT NULL
);

-- analytics-worker skriver hit, framtida dashboards läser via denna
CREATE INDEX idx_click_events_code ON click_events (code);
