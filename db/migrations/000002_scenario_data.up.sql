INSERT INTO scenario (scenario_id, cluster_id, scenario_title, scenario_type, is_default)
VALUES (1, 'Zk_default_cluster_id_for_all_scenarios', 'Exception', 'SYSTEM', true);

INSERT INTO scenario (scenario_id, cluster_id, scenario_title, scenario_type, is_default)
VALUES (2, 'Zk_default_cluster_id_for_all_scenarios', 'All Traces', 'SYSTEM', true);

INSERT INTO scenario (scenario_id, cluster_id, scenario_title, scenario_type, is_default)
VALUES (3, 'Zk_default_cluster_id_for_all_scenarios', 'Client Error', 'SYSTEM', true);

INSERT INTO scenario (scenario_id, cluster_id, scenario_title, scenario_type, is_default)
VALUES (4, 'Zk_default_cluster_id_for_all_scenarios', 'Server Error', 'SYSTEM', true);

INSERT INTO scenario (scenario_id, cluster_id, scenario_title, scenario_type, is_default)
VALUES (5, 'Zk_default_cluster_id_for_all_scenarios', 'Slow Request', 'SYSTEM', true);

INSERT INTO scenario_version (scenario_version_id, scenario_id, scenario_data, schema_version, scenario_version, created_by, created_at)
VALUES (1, 1, '{"version":"1684149787","scenario_title":"Exception","scenario_type":"SYSTEM","enabled":true,"workloads":{"55661a0e-25cb-5a1c-94cd-fad172b0caa2":{"service":"*/*","trace_role":"server","protocol":"HTTP","rule":{"type":"rule_group","condition":"AND","rules":[{"type":"rule","id":"req_method","field":"req_method","datatype":"string","input":"string","operator":"equal","value":"POST"},{"type":"rule","id":"req_path","field":"req_path","datatype":"string","input":"string","operator":"equal","value":"/exception"}]}}},"scenario_id":"1","filter":{"type":"workload","condition":"AND","workload_ids":["55661a0e-25cb-5a1c-94cd-fad172b0caa2"]},"group_by":[{"workload_id":"55661a0e-25cb-5a1c-94cd-fad172b0caa2","title":"source","hash":"source"}]}', 'v1', 1687763051, 'vaibhav', 1687763051);

INSERT INTO scenario_version (scenario_version_id, scenario_id, scenario_data, schema_version, scenario_version, created_by, created_at)
VALUES (2, 2, '{"version":"1684149744","scenario_id":"2","scenario_title":"All Traces","scenario_type":"SYSTEM","enabled":true,"workloads":{"ae427f14-9833-539c-933f-f0f26c6e4263":{"service":"*/*","trace_role":"server","protocol":"HTTP","rule":{"type":"rule_group","condition":"AND","rules":[{"type":"rule","id":"resp_status","field":"resp_status","datatype":"integer","input":"integer","operator":"greater_than","value":"199"}]}}},"filter":{"type":"workload","condition":"OR","workload_ids":["ae427f14-9833-539c-933f-f0f26c6e4263"]},"group_by":[{"workload_id":"ae427f14-9833-539c-933f-f0f26c6e4263","title":"source","hash":"source"}]}', 'v1', 1687763051, 'vaibhav', 1687763051);

INSERT INTO scenario_version (scenario_version_id, scenario_id, scenario_data, schema_version, scenario_version, created_by, created_at)
VALUES (3, 3, '{"version":"1684149744","scenario_id":"3","scenario_title":"Client Error","scenario_type":"SYSTEM","enabled":true,"workloads":{"cb567aec-13c2-573b-97d1-fb2c74756703":{"service":"*/*","trace_role":"server","protocol":"HTTP","rule":{"type":"rule_group","condition":"AND","rules":[{"type":"rule","id":"resp_status","field":"resp_status","datatype":"integer","input":"integer","operator":"greater_than","value":"399"},{"type":"rule","id":"resp_status","field":"resp_status","datatype":"integer","input":"integer","operator":"less_than","value":"500"}]}}},"filter":{"type":"workload","condition":"OR","workload_ids":["cb567aec-13c2-573b-97d1-fb2c74756703"]},"group_by":[{"workload_id":"cb567aec-13c2-573b-97d1-fb2c74756703","title":"source","hash":"source"}]}', 'v1', 1687763051, 'vaibhav', 1687763051);

INSERT INTO scenario_version (scenario_version_id, scenario_id, scenario_data, schema_version, scenario_version, created_by, created_at)
VALUES (4, 4, '{"version":"1684149744","scenario_id":"4","scenario_title":"Server Error","scenario_type":"SYSTEM","enabled":true,"workloads":{"9b9b9bd5-64dd-55bf-a11e-76027c2b0fa0":{"service":"*/*","trace_role":"server","protocol":"HTTP","rule":{"type":"rule_group","condition":"AND","rules":[{"type":"rule","id":"resp_status","field":"resp_status","datatype":"integer","input":"integer","operator":"greater_than","value":"499"}]}}},"filter":{"type":"workload","condition":"OR","workload_ids":["9b9b9bd5-64dd-55bf-a11e-76027c2b0fa0"]},"group_by":[{"workload_id":"9b9b9bd5-64dd-55bf-a11e-76027c2b0fa0","title":"destination","hash":"destination"}]}', 'v1', 1687763051, 'vaibhav', 1687763051);

INSERT INTO scenario_version (scenario_version_id, scenario_id, scenario_data, schema_version, scenario_version, created_by, created_at)
VALUES (5, 5, '{"version":"1684149744","scenario_id":"5","scenario_title":"Slow Request","scenario_type":"SYSTEM","enabled":true,"workloads":{"6d2f7f0d-94fb-5f5e-b1ba-dea30a6ddd16":{"service":"*/*","trace_role":"server","protocol":"HTTP","rule":{"type":"rule_group","condition":"AND","rules":[{"type":"rule","id":"latency","field":"latency","datatype":"integer","input":"integer","operator":"greater_than","value":"50000000"}]}}},"filter":{"type":"workload","condition":"OR","workload_ids":["6d2f7f0d-94fb-5f5e-b1ba-dea30a6ddd16"]},"group_by":[{"workload_id":"6d2f7f0d-94fb-5f5e-b1ba-dea30a6ddd16","title":"destination","hash":"destination"}]}', 'v1', 1687763051, 'vaibhav', 1687763051);