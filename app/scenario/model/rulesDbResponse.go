package model

type ScenarioDbResponse struct {
	ClusterId       string `json:"cluster_id"`
	ScenarioData    string `json:"scenario_data"`
	ScenarioTitle   string `json:"scenario_title"`
	ScenarioType    string `json:"scenario_type"`
	IsDefault       bool   `json:"is_default"`
	ScenarioId      string `json:"scenario_id"`
	SchemaVersion   string `json:"schema_version"`
	ScenarioVersion int64  `json:"scenario_version"`
	Deleted         bool   `json:"deleted"`
	DeletedBy       string `json:"deleted_by"`
	DeletedAt       int64  `json:"deleted_at"`
	CreatedBy       string `json:"created_by"`
	CreatedAt       int64  `json:"created_at"`
	Disabled        bool   `json:"disabled"`
	DisabledBy      string `json:"disabled_by"`
	DisabledAt      *int64 `json:"disabled_at"`
	UpdatedAt       int64  `json:"updated_at"`
}
