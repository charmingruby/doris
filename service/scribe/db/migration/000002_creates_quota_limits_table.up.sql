CREATE TABLE quota_limits (
    id VARCHAR PRIMARY KEY,
    quota_id VARCHAR NOT NULL,
    kind TEXT NOT NULL,
    max_value INTEGER NOT NULL,
    unit TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_quota_id ON quotas (id);