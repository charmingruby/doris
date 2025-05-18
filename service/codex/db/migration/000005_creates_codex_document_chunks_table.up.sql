CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE codex_document_chunks (
    id SERIAL PRIMARY KEY,
    codex_document_id VARCHAR NOT NULL,
    excerpt VARCHAR,
    embedding VECTOR(1024), -- mistral dimension
    created_at TIMESTAMP NOT NULL,

    FOREIGN KEY (codex_document_id) REFERENCES codex_documents(id)
);
