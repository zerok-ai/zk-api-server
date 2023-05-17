package model

type DataTypes string
type InputTypes = string
type OperatorTypes = string
type ValueTypes interface {
}
type Rule struct {
	ID       string        `json:"id,omitempty"`
	Field    string        `json:"field,omitempty"`
	Type     DataTypes     `json:"type,omitempty"`
	Input    InputTypes    `json:"input,omitempty"`
	Operator OperatorTypes `json:"operator,omitempty"`
	Key      string        `json:"key,omitempty"`
	Value    ValueTypes    `json:"value,omitempty"`
}

type FilterRule struct {
	Condition *string      `json:"condition,omitempty"`
	Service   *string      `json:"service,omitempty"`
	TraceRole *string      `json:"trace_role,omitempty"`
	Rules     []FilterType `json:"rules,omitempty"`
}

type FilterType struct {
	FilterRule
	Rule
}

type NewRuleSchema struct {
	Version   int                   `json:"version"`
	Workloads map[string]FilterType `json:"workloads"`
	FilterId  string                `json:"filter_id"`
	Filters   Filters               `json:"filters"`
}
type Filters struct {
	Type      string   `json:"type"`
	Condition string   `json:"condition"`
	Workloads []string `json:"workloads"`
}
