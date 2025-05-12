CREATE TABLE quotas (
    id VARCHAR PRIMARY KEY,
    tier VARCHAR NOT NULL,
    kind VARCHAR NOT NULL,
    max_value INT NOT NULL,
    unit VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_quotas_tier_kind ON quotas (tier, kind);