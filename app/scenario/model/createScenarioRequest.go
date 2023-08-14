package model

type RuleGroup struct {
	Type      string `json:"type"`
	Condition string `json:"condition"`
	Rules     []Rule `json:"rules"`
}

type Rule struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	Field    string `json:"field"`
	Datatype string `json:"datatype"`
	Input    string `json:"input"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Workload struct {
	Service   string    `json:"service"`
	TraceRole string    `json:"trace_role"`
	Protocol  string    `json:"protocol"`
	Rule      RuleGroup `json:"rule"`
}

type GroupByItem struct {
	WorkloadIndex int    `json:"workload_index"`
	Title         string `json:"title"`
	Hash          string `json:"hash"`
}

type CreateScenarioRequest struct {
	ScenarioTitle string        `json:"scenario_title"`
	ScenarioType  string        `json:"scenario_type"`
	Workloads     []Workload    `json:"workloads"`
	GroupBy       []GroupByItem `json:"group_by"`
}
