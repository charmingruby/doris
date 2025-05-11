CREATE TABLE quota_limit_usages (
    id VARCHAR PRIMARY KEY,
    correlation_id VARCHAR NOT NULL,
    quota_limit_id VARCHAR NOT NULL,
    current_usage INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_reset_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_quota_limit_usages_correlation_id ON quota_limit_usages (correlation_id);
CREATE INDEX IF NOT EXISTS idx_quota_limit_usages_quota_limit_id ON quota_limits (id);