package dto

import (
	"encoding/json"
	"time"
)

type Level string
type Type string

const (
	Org        Level = "ORG"
	Cluster    Level = "CLUSTER"
	Prometheus Type  = "PROMETHEUS"
)

type Integration struct {
	ID             *string         `json:"id"`
	ClusterId      string          `json:"cluster_id"`
	Alias          string          `json:"alias"`
	Type           Type            `json:"type"`
	URL            string          `json:"url"`
	Authentication json.RawMessage `json:"authentication"`
	Level          Level           `json:"level"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Deleted        bool            `json:"deleted"`
	Disabled       bool            `json:"disabled"`
}

func (integration Integration) GetAllColumns() []any {
	return []any{integration.ClusterId, integration.Alias, integration.Type, integration.URL, integration.Authentication, integration.Level, integration.CreatedAt, integration.UpdatedAt, integration.Deleted, integration.Disabled}
}

type IntegrationRequest struct {
	ID             *string `json:"id"`
	ClusterId      string
	Alias          string          `json:"alias"`
	Type           Type            `json:"type"`
	URL            string          `json:"url"`
	Authentication json.RawMessage `json:"authentication"`
	Level          Level           `json:"level"`
	Deleted        bool            `json:"deleted"`
	Disabled       bool            `json:"disabled"`
}
