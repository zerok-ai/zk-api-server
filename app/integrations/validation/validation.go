package validation

import (
	"errors"
	"github.com/zerok-ai/zk-utils-go/common"
	"zk-api-server/app/integrations/model/dto"
)

func ValidateIntegrationsUpsertRequest(integration dto.IntegrationRequest) error {

	switch integration.Level {
	case dto.Org, dto.Cluster:
	default:
		return errors.New("invalid 'level' value. Allowed values are 'org' or 'cluster'")
	}

	if integration.Type != dto.Prometheus {
		return errors.New("invalid 'type' value. Allowed value is 'PROMETHEUS'")
	}

	if common.IsEmpty(integration.ClusterId) {
		return errors.New("clusterId cannot be empty")
	}

	if common.IsEmpty(integration.Alias) {
		return errors.New("alias cannot be empty")
	}

	return nil
}
