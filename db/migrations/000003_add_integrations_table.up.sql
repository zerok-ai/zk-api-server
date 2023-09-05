CREATE TABLE IF NOT EXISTS integrations
(
    id             SERIAL PRIMARY KEY,
    type           VARCHAR(255),
    url            VARCHAR(255),
    authentication JSONB,
    level          VARCHAR(30),
    created_at     TIMESTAMP,
    updated_at     TIMESTAMP,
    deleted        BOOLEAN DEFAULT FALSE,
    disabled       BOOLEAN DEFAULT FALSE
);
