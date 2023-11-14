package service

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"net/http"
	"zk-api-server/app/integrations/model/dto"
	"zk-api-server/app/integrations/model/transformer"
	"zk-api-server/app/integrations/repository"
	"zk-api-server/app/utils/errors"
)

type IntegrationsService interface {
	GetAllIntegrations(clusterId string, onlyActive bool) (transformer.IntegrationResponse, *zkerrors.ZkError)
	UpsertIntegration(integration dto.Integration) (bool, *string, *zkerrors.ZkError)
	GetAnIntegrationDetails(insertId string) (int, *zkerrors.ZkError)
	GetIntegrationsStatusResponse(clusterId, url, username, password string) (*http.Response, *zkerrors.ZkError)
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

func (i integrationsService) GetAnIntegrationDetails(insertId string) (int, *zkerrors.ZkError) {
	integration, err := i.repo.GetAnIntegrationDetails(insertId)
	if err != nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return iris.StatusInternalServerError, &zkError
	} else if integration == nil || len(integration) == 0 {
		zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestInvalidClusterAndUrlCombination, err)
		return iris.StatusBadRequest, &zkError
	}

	var auth dto.Auth
	var username, password string
	err = json.Unmarshal(integration[0].Authentication, &auth)
	if err == nil {
		username = auth.Username
		password = auth.Password
	}

	httpResp, zkErr := i.GetIntegrationsStatusResponse(integration[0].ClusterId, integration[0].URL, username, password)
	if zkErr != nil {
		zkLogger.Error(LogTag, "Error while getting the integration status: ", zkErr)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, zkErr)
		return iris.StatusInternalServerError, &newZkErr
	}

	return httpResp.StatusCode, nil
}

func (i integrationsService) UpsertIntegration(integration dto.Integration) (bool, *string, *zkerrors.ZkError) {
	if integration.ID != nil {
		if row, err := i.repo.GetIntegrationsById(*integration.ID, integration.ClusterId); err != nil || row == nil {
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return false, nil, &zkError
		} else if row != nil {
			if valid := validateIntegrationsForUpsert(*row, integration); !valid {
				zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, nil)
				zkLogger.Error(LogTag, "Integration validation failed")
				return false, nil, &zkError
			}
		}

		done, err := i.repo.UpdateIntegration(integration)
		if err != nil {
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return false, nil, &zkError
		}

		return done, integration.ID, nil
	}

	done, id, err := i.repo.InsertIntegration(integration)
	if err != nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return false, nil, &zkError
	}

	return done, common.ToString(id), nil
}

func (i integrationsService) GetIntegrationsStatusResponse(clusterId, url, username, password string) (*http.Response, *zkerrors.ZkError) {
	if common.IsEmpty(clusterId) || common.IsEmpty(url) {
		zkLogger.Error(LogTag, "ClusterId or url is empty")
		zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestInvalidClusterAndUrlCombination, nil)
		return nil, &zkError
	}

	prometheusStatusQueryPath := "/api/v1/query?query=up"
	return zkHttp.Create().
		BasicAuth(username, password).
		Header("X-PROXY-DESTINATION", url+prometheusStatusQueryPath).
		Header("X-CLIENT-ID", clusterId).
		Get("http://zk-wsp-server.zkcloud.svc.cluster.local:8989/request")
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
