package dto

import (
	"encoding/json"
	"time"

	zkIntegration "github.com/zerok-ai/zk-utils-go/integration/model"
)

type Integration struct {
	ID             *string             `json:"id"`
	ClusterId      string              `json:"cluster_id"`
	Alias          string              `json:"alias"`
	Type           zkIntegration.Type  `json:"type"`
	URL            string              `json:"url"`
	Authentication json.RawMessage     `json:"authentication"`
	Level          zkIntegration.Level `json:"level"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	Deleted        bool                `json:"deleted"`
	Disabled       bool                `json:"disabled"`
	MetricServer   bool                `json:"metric_server"`
}

func (integration Integration) GetAllColumns() []any {
	return []any{integration.ClusterId, integration.Alias, integration.Type, integration.URL, integration.Authentication, integration.Level, integration.CreatedAt, integration.UpdatedAt, integration.Deleted, integration.Disabled, integration.MetricServer}
}

type IntegrationRequest struct {
	ID             *string `json:"id"`
	ClusterId      string
	Alias          string              `json:"alias"`
	Type           zkIntegration.Type  `json:"type"`
	URL            string              `json:"url"`
	Authentication json.RawMessage     `json:"authentication"`
	Level          zkIntegration.Level `json:"level"`
	Deleted        bool                `json:"deleted"`
	Disabled       bool                `json:"disabled"`
	MetricServer   bool                `json:"metric_server"`
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpsertIntegrationResponse struct {
	IntegrationId string `json:"integration_id"`
	Status        int    `json:"status"`
}
