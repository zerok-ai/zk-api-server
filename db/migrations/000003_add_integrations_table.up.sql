CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS integrations
(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    cluster_id     VARCHAR(255) NOT NULL,
    alias         VARCHAR(255) NOT NULL,
    type           VARCHAR(255) NOT NULL,
    url            VARCHAR(255),
    authentication JSONB,
    level          VARCHAR(30) NOT NULL,
    created_at     TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP NOT NULL,
    deleted        BOOLEAN DEFAULT FALSE,
    disabled       BOOLEAN DEFAULT FALSE,
    metric_server BOOLEAN DEFAULT FALSE
);