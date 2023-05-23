package model

type RulesDbResponse struct {
	ClusterId     string `json:"cluster_id"`
	Filters       string `json:"filters"`
	FilterId      string `json:"filter_id"`
	SchemaVersion string `json:"schema_version"`
	Version       int64  `json:"version"`
	UpdatedAt     int64  `json:"updated_at"`
	Deleted       bool   `json:"deleted"`
	DeletedBy     string `json:"deleted_by"`
	DeletedAt     int64  `json:"deleted_at"`
	CreatedBy     string `json:"created_by"`
	CreatedAt     int64  `json:"created_at"`
	IsDeleted     bool   `json:"is_deleted"`
}
