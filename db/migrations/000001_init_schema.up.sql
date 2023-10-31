CREATE TABLE IF NOT EXISTS zk_scenario
(
    scenario_id    SERIAL PRIMARY KEY,
    cluster_id     VARCHAR(255),
    scenario_title VARCHAR(255),
    disabled       BOOL         DEFAULT FALSE,
    disabled_by    VARCHAR(255) DEFAULT NULL,
    disabled_at    BIGINT       DEFAULT NULL,
    scenario_type  VARCHAR(50),
    is_default     BOOLEAN      DEFAULT false,
    deleted        BOOLEAN      DEFAULT FALSE,
    deleted_by     VARCHAR(255) DEFAULT NULL,
    deleted_at     BIGINT       DEFAULT NULL,
    updated_at     BIGINT       NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())
);

ALTER SEQUENCE zk_scenario_scenario_id_seq RESTART WITH 1000;

CREATE TABLE IF NOT EXISTS zk_scenario_version
(
    scenario_version_id SERIAL PRIMARY KEY,
    scenario_id         INTEGER REFERENCES zk_scenario (scenario_id) ON DELETE CASCADE,
    scenario_data       BYTEA,
    schema_version      VARCHAR(255),
    scenario_version    BIGINT NOT NULL NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    created_by          VARCHAR(255),
    created_at          BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())
);

ALTER SEQUENCE zk_scenario_version_scenario_version_id_seq RESTART WITH 1000;

CREATE OR REPLACE FUNCTION zk_update_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = EXTRACT(EPOCH FROM NOW());
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER zk_update_scenarios_updated_at
    BEFORE UPDATE ON zk_scenario
    FOR EACH ROW
EXECUTE FUNCTION zk_update_updated_at();