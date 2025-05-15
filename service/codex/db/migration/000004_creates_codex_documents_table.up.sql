CREATE TABLE codex_documents (
    id VARCHAR PRIMARY KEY,
    codex_id VARCHAR NOT NULL,
    title VARCHAR NOT NULL,
    image_url VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,

    FOREIGN KEY (codex_id) REFERENCES codex(id)
);