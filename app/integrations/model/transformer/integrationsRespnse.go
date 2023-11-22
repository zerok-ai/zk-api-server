package transformer

import (
	"time"
	"zk-api-server/app/integrations/model/dto"

	zkIntegrationResponse "github.com/zerok-ai/zk-utils-go/integration/model"
)

type IntegrationResponse struct {
	Response []zkIntegrationResponse.IntegrationResponseObj `json:"integrations"`
}

func FromIntegrationArrayToIntegrationResponse(iArr []dto.Integration) IntegrationResponse {
	responseArr := make([]zkIntegrationResponse.IntegrationResponseObj, 0)
	for _, i := range iArr {
		responseArr = append(responseArr, zkIntegrationResponse.IntegrationResponseObj{
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
			MetricServer:   i.MetricServer,
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
		MetricServer:   iReq.MetricServer,
	}
}
