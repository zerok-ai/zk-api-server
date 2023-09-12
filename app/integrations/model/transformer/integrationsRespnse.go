package transformer

import (
	"encoding/json"
	"time"
	"zk-api-server/app/integrations/model/dto"
)

type IntegrationResponseObj struct {
	ID             string          `json:"id"`
	ClusterId      string          `json:"cluster_id,omitempty"`
	Alias          string          `json:"alias"`
	Type           dto.Type        `json:"type"`
	URL            string          `json:"url"`
	Authentication json.RawMessage `json:"authentication"`
	Level          dto.Level       `json:"level"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Deleted        bool            `json:"deleted"`
	Disabled       bool            `json:"disabled"`
}

type IntegrationResponse struct {
	Response []IntegrationResponseObj `json:"integrations"`
}

func FromIntegrationArrayToIntegrationResponse(iArr []dto.Integration) IntegrationResponse {
	responseArr := make([]IntegrationResponseObj, 0)
	for _, i := range iArr {
		responseArr = append(responseArr, IntegrationResponseObj{
			ID:             *i.ID,
			ClusterId:      i.ClusterId,
			Alias:          i.Alias,
			Type:           i.Type,
			URL:            i.URL,
			Authentication: i.Authentication,
			Level:          i.Level,
			CreatedAt:      i.CreatedAt,
			UpdatedAt:      i.UpdatedAt,
			Deleted:        i.Deleted,
			Disabled:       i.Disabled,
		})
	}

	return IntegrationResponse{Response: responseArr}
}

func FromIntegrationsRequestToIntegrationsDto(iReq dto.IntegrationRequest) dto.Integration {
	currentTime := time.Now()
	return dto.Integration{
		ID:             iReq.ID,
		ClusterId:      iReq.ClusterId,
		Alias:          iReq.Alias,
		Type:           iReq.Type,
		URL:            iReq.URL,
		Authentication: iReq.Authentication,
		Level:          iReq.Level,
		CreatedAt:      currentTime,
		UpdatedAt:      currentTime,
		Deleted:        iReq.Deleted,
		Disabled:       iReq.Disabled,
	}
}
