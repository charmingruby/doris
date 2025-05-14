CREATE TABLE quota_usages (
    id VARCHAR PRIMARY KEY,
    correlation_id VARCHAR NOT NULL,
    quota_id VARCHAR NOT NULL,
    current_usage INT NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_reset_at TIMESTAMP,
    updated_at TIMESTAMP,

    CONSTRAINT fk_quota_usages_quota_id FOREIGN KEY (quota_id) REFERENCES quotas(id)
);

CREATE INDEX idx_quota_usages_quota_id ON quota_usages (quota_id);
CREATE INDEX idx_quota_usages_correlation_id ON quota_usages (correlation_id);
