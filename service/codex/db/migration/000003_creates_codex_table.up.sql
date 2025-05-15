CREATE TABLE codex (
    id VARCHAR PRIMARY KEY,
    correlation_id VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);