drop table zk_scenario_version;
drop table zk_scenario;
DROP TRIGGER IF EXISTS zk_update_scenarios_updated_at ON zk_scenario;
DROP FUNCTION IF EXISTS zk_update_updated_at();