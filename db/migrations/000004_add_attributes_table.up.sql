CREATE TABLE IF NOT EXISTS otel_attributes
(
    id             SERIAL PRIMARY KEY,
    version        VARCHAR(255) NOT NULL,
    key_set        VARCHAR(255) NOT NULL,
    attribute_list JSONB NOT NULL,
    CONSTRAINT unique_version_key_set UNIQUE (version, key_set)
);
