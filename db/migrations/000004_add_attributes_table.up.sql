CREATE TABLE IF NOT EXISTS zk_attributes
(
    id             SERIAL PRIMARY KEY,
    version        VARCHAR(255) NOT NULL,
    protocol      VARCHAR(255) NOT NULL,
    executor      VARCHAR(255) NOT NULL,
    attribute_list JSONB        NOT NULL,
    updated_at     BIGINT       NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    CONSTRAINT unique_version_key_set UNIQUE (version, protocol, executor)
);

CREATE TRIGGER zk_update_attributes_updated_at
    BEFORE UPDATE ON zk_attributes
    FOR EACH ROW
EXECUTE FUNCTION zk_update_updated_at();
