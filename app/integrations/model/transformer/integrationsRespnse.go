package transformer

import (
	"encoding/json"
	"time"
	"zk-api-server/app/integrations/model/dto"
)

type IntegrationResponseObj struct {
	ID             int             `json:"id"`
	ClusterId      string          `json:"cluster_id"`
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
			ID:             i.ID,
			ClusterId:      i.ClusterId,
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
	return dto.Integration{
		ID:             iReq.ID,
		Type:           iReq.Type,
		URL:            iReq.URL,
		Authentication: iReq.Authentication,
		Level:          iReq.Level,
		CreatedAt:      iReq.CreatedAt,
		UpdatedAt:      iReq.UpdatedAt,
		Deleted:        iReq.Deleted,
		Disabled:       iReq.Disabled,
	}
}
