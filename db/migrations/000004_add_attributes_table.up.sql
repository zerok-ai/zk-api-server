CREATE TABLE IF NOT EXISTS attributes
(
    id             SERIAL PRIMARY KEY,
    version        VARCHAR(255) NOT NULL,
    key_set        VARCHAR(255) NOT NULL,
    protocol      VARCHAR(255) NOT NULL,
    executor      VARCHAR(255) NOT NULL,
    attribute_list JSONB        NOT NULL,
    updated_at     BIGINT       NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    CONSTRAINT unique_version_key_set UNIQUE (version, key_set)
);

CREATE TABLE IF NOT EXISTS protocol_attributes_mapping
(
    id                 SERIAL PRIMARY KEY,
    protocol           VARCHAR(255) NOT NULL,
    otel_attributes_id INTEGER      NOT NULL REFERENCES attributes (id)
);
