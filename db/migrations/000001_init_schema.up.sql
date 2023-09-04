CREATE TABLE IF NOT EXISTS scenario
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
    updated_at     BIGINT       NOT NULL NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())
);

ALTER SEQUENCE scenario_scenario_id_seq RESTART WITH 1000;

CREATE TABLE IF NOT EXISTS scenario_version
(
    scenario_version_id SERIAL PRIMARY KEY,
    scenario_id         INTEGER REFERENCES scenario (scenario_id) ON DELETE CASCADE,
    scenario_data       BYTEA,
    schema_version      VARCHAR(255),
    scenario_version    BIGINT NOT NULL NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    created_by          VARCHAR(255),
    created_at          BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())
);

ALTER SEQUENCE scenario_version_scenario_version_id_seq RESTART WITH 1000;