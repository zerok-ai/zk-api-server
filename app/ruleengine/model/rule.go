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
	Condition     *string      `json:"condition,omitempty"`
	ZkRequestType *Rule        `json:"zk_request_type,omitempty"`
	Rules         []FilterType `json:"rules,omitempty"`
}

type FilterType struct {
	FilterRule
	Rule
}

type Value interface {
	string | int | float64 | SourceDestinationHolder
}

type SourceDestinationHolder struct {
	ServiceName string `json:"service_name"`
	Ip          string `json:"ip"`
	PodName     string `json:"pod_name"`
}
