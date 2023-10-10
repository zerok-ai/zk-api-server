INSERT INTO zk_scenario (scenario_id, cluster_id, scenario_title, scenario_type, is_default)
VALUES (1, 'Zk_default_cluster_id_for_all_scenarios', 'Exception', 'SYSTEM', true);

INSERT INTO zk_scenario (scenario_id, cluster_id, scenario_title, scenario_type, is_default)
VALUES (2, 'Zk_default_cluster_id_for_all_scenarios', '4xx Error', 'SYSTEM', true);

INSERT INTO zk_scenario (scenario_id, cluster_id, scenario_title, scenario_type, is_default)
VALUES (3, 'Zk_default_cluster_id_for_all_scenarios', '5xx Error', 'SYSTEM', true);

INSERT INTO zk_scenario (scenario_id, cluster_id, scenario_title, scenario_type, is_default)
VALUES (4, 'Zk_default_cluster_id_for_all_scenarios', 'Requests with latency > 100ms', 'SYSTEM', true);

INSERT INTO zk_scenario_version (scenario_version_id, scenario_id, scenario_data, scenario_version, schema_version, created_by)
VALUES (1, 1, '{"version":"1684149787","scenario_title":"Exception","scenario_type":"SYSTEM","enabled":true,"workloads":{"6e1a43c1-bf8e-535b-90dc-01bb475fa4f7":{"executor":"OTEL","service":"*/*","trace_role":"server","protocol":"HTTP","rule":{"type":"rule_group","condition":"AND","rules":[{"type":"rule","id":"http_response_status_code","field":"Response Status","datatype":"","input":"","operator":"exists","value":""},{"type":"rule","id":"errors","field":"Errors","datatype":"","input":"","operator":"exists","value":""}]}}},"scenario_id":"1","filter":{"type":"workload","condition":"AND","workload_ids":["6e1a43c1-bf8e-535b-90dc-01bb475fa4f7"]},"group_by":[{"workload_id":"6e1a43c1-bf8e-535b-90dc-01bb475fa4f7","title":"source_service","hash":"source_service"},{"workload_id":"6e1a43c1-bf8e-535b-90dc-01bb475fa4f7","title":"errors_message","hash":"errors_message"}],"rate_limit":[{"bucket_max_size":5,"bucket_refill_size":5,"tick_duration":"1m"}]}', 1693843788, 'v1', 'SYSTEM');

INSERT INTO zk_scenario_version (scenario_version_id, scenario_id, scenario_data, scenario_version, schema_version, created_by)
VALUES (2, 2, '{"version":"1684149744","scenario_id":"2","scenario_title":"4xx Error","scenario_type":"SYSTEM","enabled":true,"workloads":{"db091e37-e1c8-5c75-b257-3980e7e49371":{"executor":"OTEL","service":"*/*","trace_role":"server","protocol":"HTTP","rule":{"type":"rule_group","condition":"AND","rules":[{"type":"rule","id":"http_response_status_code","field":"Response Status","datatype":"integer","input":"integer","operator":"greater_than","value":"399"},{"type":"rule","id":"http_response_status_code","field":"Response Status","datatype":"integer","input":"integer","operator":"less_than","value":"500"}]}}},"filter":{"type":"workload","condition":"AND","workload_ids":["db091e37-e1c8-5c75-b257-3980e7e49371"]},"group_by":[{"workload_id":"db091e37-e1c8-5c75-b257-3980e7e49371","title":"source_service","hash":"source_service"},{"workload_id":"db091e37-e1c8-5c75-b257-3980e7e49371","title":"http_response_status_code","hash":"http_response_status_code"}],"rate_limit":[{"bucket_max_size":5,"bucket_refill_size":5,"tick_duration":"1m"}]}', 1693843788, 'v1', 'SYSTEM');

INSERT INTO zk_scenario_version (scenario_version_id, scenario_id, scenario_data, scenario_version, schema_version, created_by)
VALUES (3, 3, '{"version":"1684149744","scenario_id":"3","scenario_title":"5xx Error","scenario_type":"SYSTEM","enabled":true,"workloads":{"cc13b872-b3d9-52ed-aa64-210114c0cbef":{"executor":"OTEL","service":"*/*","trace_role":"server","protocol":"HTTP","rule":{"type":"rule_group","condition":"AND","rules":[{"type":"rule","id":"http_response_status_code","field":"Status Code","datatype":"integer","input":"integer","operator":"greater_than","value":"499"}]}}},"filter":{"type":"workload","condition":"AND","workload_ids":["cc13b872-b3d9-52ed-aa64-210114c0cbef"]},"group_by":[{"workload_id":"cc13b872-b3d9-52ed-aa64-210114c0cbef","title":"dest_service","hash":"dest_service"},{"workload_id":"cc13b872-b3d9-52ed-aa64-210114c0cbef","title":"http_response_status_code","hash":"http_response_status_code"}],"rate_limit":[{"bucket_max_size":5,"bucket_refill_size":5,"tick_duration":"1m"}]}', 1693843788, 'v1', 'SYSTEM');

INSERT INTO zk_scenario_version (scenario_version_id, scenario_id, scenario_data, scenario_version, schema_version, created_by)
VALUES (4, 4, '{"version":"1684149744","scenario_id":"4","scenario_title":"Requests with latency > 100ms","scenario_type":"SYSTEM","enabled":true,"workloads":{"4dfca142-5bf8-5762-94ac-9bf5300f985c":{"executor":"OTEL","service":"*/*","trace_role":"server","protocol":"HTTP","rule":{"type":"rule_group","condition":"AND","rules":[{"type":"rule","id":"latency_ns","field":"Latency in ns","datatype":"integer","input":"integer","operator":"greater_than","value":"100000000"}]}}},"filter":{"type":"workload","condition":"AND","workload_ids":["4dfca142-5bf8-5762-94ac-9bf5300f985c"]},"group_by":[{"workload_id":"4dfca142-5bf8-5762-94ac-9bf5300f985c","title":"dest_service","hash":"dest_service"}],"rate_limit":[{"bucket_max_size":5,"bucket_refill_size":5,"tick_duration":"1m"}]}', 1693843788, 'v1', 'SYSTEM');
