CREATE TABLE IF NOT EXISTS api_keys
(
    id varchar PRIMARY KEY NOT NULL,
    first_name varchar NOT NULL,
    last_name varchar NOT NULL,
    email varchar NOT NULL,
    key varchar NOT NULL,
    tier varchar NOT NULL,
    status varchar NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_key ON api_keys (key);
CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_email ON api_keys (email);
