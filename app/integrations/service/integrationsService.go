package service

import (
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"zk-api-server/app/integrations/model/dto"
	"zk-api-server/app/integrations/model/transformer"
	"zk-api-server/app/integrations/repository"
)

type IntegrationsService interface {
	GetAllIntegrations(clusterId string, onlyActive bool) (transformer.IntegrationResponse, *zkerrors.ZkError)
	UpsertIntegration(integration dto.Integration) (bool, *zkerrors.ZkError)
}

var LogTag = "integrations_service"

type integrationsService struct {
	repo repository.IntegrationRepo
}

func NewIntegrationsService(repo repository.IntegrationRepo) IntegrationsService {
	return &integrationsService{repo: repo}
}

func (i integrationsService) GetAllIntegrations(clusterId string, onlyActive bool) (transformer.IntegrationResponse, *zkerrors.ZkError) {
	integrations, err := i.repo.GetAllIntegrations(clusterId, onlyActive)
	if err != nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return transformer.IntegrationResponse{}, &zkError
	}

	return transformer.FromIntegrationArrayToIntegrationResponse(integrations), nil
}

func (i integrationsService) UpsertIntegration(integration dto.Integration) (bool, *zkerrors.ZkError) {
	if integration.ID != nil {
		if row, err := i.repo.GetIntegrationsById(*integration.ID, integration.ClusterId); err != nil || row == nil {
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return false, &zkError
		} else if row != nil {
			if valid := validateIntegrationsForUpsert(*row, integration); !valid {
				zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, nil)
				zkLogger.Error(LogTag, "Integration validation failed")
				return false, &zkError
			}
		}

		done, err := i.repo.UpdateIntegration(integration)
		if err != nil {
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return false, &zkError
		}

		return done, nil
	}

	done, err := i.repo.InsertIntegration(integration)
	if err != nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return false, &zkError
	}

	return done, nil
}

func validateIntegrationsForUpsert(fromDb, fromRequest dto.Integration) bool {
	if *fromDb.ID != *fromRequest.ID {
		zkLogger.Error(LogTag, "Integration validation failed different id")
		return false
	}

	if fromDb.ClusterId != fromRequest.ClusterId {
		zkLogger.Error(LogTag, "Integration validation failed different clusterId")
		return false
	}

	if fromDb.Type != fromRequest.Type {
		zkLogger.Error(LogTag, "Integration validation failed different type")
		return false
	}

	return true
}
