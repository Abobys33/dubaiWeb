CREATE TABLE IF NOT EXISTS feedback_requests
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    phone      TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_feedback_requests_created_at
    ON feedback_requests (created_at);