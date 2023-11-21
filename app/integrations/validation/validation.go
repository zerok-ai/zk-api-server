package validation

import (
	"errors"
	"github.com/zerok-ai/zk-utils-go/common"
	zkIntegration "github.com/zerok-ai/zk-utils-go/integration/model"
	"zk-api-server/app/integrations/model/dto"
)

func ValidateIntegrationsUpsertRequest(integration dto.IntegrationRequest) error {

	switch integration.Level {
	case zkIntegration.Org, zkIntegration.Cluster:
	default:
		return errors.New("invalid 'level' value. Allowed values are 'ORG' or 'CLUSTER'")
	}

	if integration.Type != zkIntegration.Prometheus {
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
