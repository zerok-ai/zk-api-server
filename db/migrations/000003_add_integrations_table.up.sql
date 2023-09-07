CREATE TABLE IF NOT EXISTS integrations
(
    id             SERIAL PRIMARY KEY,
    cluster_id     VARCHAR(255) NOT NULL,
    type           VARCHAR(255) NOT NULL,
    url            VARCHAR(255),
    authentication JSONB,
    level          VARCHAR(30) NOT NULL,
    created_at     TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP NOT NULL,
    deleted        BOOLEAN DEFAULT FALSE,
    disabled       BOOLEAN DEFAULT FALSE
);

ALTER SEQUENCE integrations_id_seq RESTART WITH 1;
