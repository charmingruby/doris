CREATE TABLE IF NOT EXISTS otps 
(
    id varchar PRIMARY KEY NOT NULL,
    code varchar NOT NULL,
    purpose varchar NOT NULL,
    correlation_id varchar NOT NULL,
    expires_at timestamp NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_otps_correlation_id ON otps (correlation_id);