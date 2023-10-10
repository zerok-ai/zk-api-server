CREATE TABLE IF NOT EXISTS attributes
(
    id             SERIAL PRIMARY KEY,
    version        VARCHAR(255) NOT NULL,
    protocol      VARCHAR(255) NOT NULL,
    executor      VARCHAR(255) NOT NULL,
    attribute_list JSONB        NOT NULL,
    updated_at     BIGINT       NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    CONSTRAINT unique_version_key_set UNIQUE (version, protocol, executor)
);

CREATE OR REPLACE FUNCTION update_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = EXTRACT(EPOCH FROM NOW());
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_attributes_updated_at
    BEFORE UPDATE ON attributes
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_scenarios_updated_at
    BEFORE UPDATE ON scenario
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
