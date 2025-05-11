CREATE TABLE quotas (
    id VARCHAR PRIMARY KEY,
    tier VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);