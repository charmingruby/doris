CREATE TABLE qas (
    id VARCHAR PRIMARY KEY,
    codex_id VARCHAR NOT NULL,
    question VARCHAR NOT NULL,
    answer VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,

    FOREIGN KEY (codex_id) REFERENCES codex(id)
);