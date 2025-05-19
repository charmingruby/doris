CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE codex_document_chunks (
    id VARCHAR PRIMARY KEY,
    codex_document_id VARCHAR NOT NULL,
    embedding VECTOR(768), -- tinyllm dimension
    created_at TIMESTAMP NOT NULL,

    FOREIGN KEY (codex_document_id) REFERENCES codex_documents(id)
);
