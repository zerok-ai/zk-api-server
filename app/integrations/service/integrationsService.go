package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"io"
	"net/http"
	"zk-api-server/app/integrations/model/dto"
	"zk-api-server/app/integrations/model/transformer"
	"zk-api-server/app/integrations/repository"
	"zk-api-server/app/utils"
	"zk-api-server/app/utils/errors"
)

type IntegrationsService interface {
	GetAllIntegrations(clusterId string, onlyActive bool) (transformer.IntegrationResponse, *zkerrors.ZkError)
	UpsertIntegration(integration dto.Integration) (bool, *string, *zkerrors.ZkError)
	TestIntegrationConnection(integrationId string) (dto.TestConnectionResponse, *zkerrors.ZkError)
	TestUnSyncedIntegrationConnection(integration dto.Integration) (dto.TestConnectionResponse, *zkerrors.ZkError)
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

func (i integrationsService) TestIntegrationConnection(integrationId string) (dto.TestConnectionResponse, *zkerrors.ZkError) {
	var resp dto.TestConnectionResponse
	integration, zkError := getIntegrationDetails(i, integrationId)
	if zkError != nil {
		return resp, zkError
	}

	httpResp, zkErr := getPrometheusApiResponse(integration[0])
	if zkErr != nil {
		zkLogger.Error(LogTag, "Error while getting the integration status: ", zkErr)
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, zkErr)
		return resp, &zkErr
	}

	if httpResp.StatusCode != iris.StatusOK {
		zkLogger.Error(LogTag, "Status Code not 200")
		resp.ConnectionStatus = utils.StatusError
		resp.ConnectionMessage = httpResp.Status
		return resp, nil
	} else {
		respBody, err := io.ReadAll(httpResp.Body)
		if err != nil {
			zkLogger.Error(LogTag, "Error while reading the response body: ", err)
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return resp, &zkError
		}

		body := map[string]dto.TestConnectionResponse{}
		err = json.Unmarshal(respBody, &body)
		if err != nil {
			zkLogger.Error(LogTag, "Error while unmarshalling the response body: ", err)
			newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return resp, &newZkErr
		}

		return body["payload"], nil
	}
}

func (i integrationsService) TestUnSyncedIntegrationConnection(integration dto.Integration) (dto.TestConnectionResponse, *zkerrors.ZkError) {
	var resp dto.TestConnectionResponse
	if common.IsEmpty(integration.URL) {
		zkLogger.Error(LogTag, "url is empty")
		zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestInvalidClusterAndUrlCombination, nil)
		return resp, &zkError
	}

	username, password := getUsernamePassword(integration)
	body := struct {
		Url string `json:"url"`
		dto.Auth
	}{
		Url: integration.URL,
		Auth: dto.Auth{
			Username: username,
			Password: password,
		},
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		zkLogger.Error(LogTag, "Error while marshalling the request body: ", err)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return resp, &newZkErr
	}
	reader := bytes.NewReader(reqBody)

	response, zkErr := zkHttp.
		Create().
		Header("X-PROXY-DESTINATION", "http://zk-axon.zk-client.svc.cluster.local:80/v1/c/axon/prom/unsaved/status").
		Header("X-CLIENT-ID", integration.ClusterId).
		Post("http://zk-wsp-server.zkcloud.svc.cluster.local:8989/request", reader)

	if zkErr != nil {
		zkLogger.Error(LogTag, "Error while getting the integration status: ", zkErr)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, zkErr)
		return resp, &newZkErr
	}

	if response.StatusCode != iris.StatusOK {
		zkLogger.Error(LogTag, "Status Code not 200")
		resp.ConnectionStatus = utils.StatusError
		resp.ConnectionMessage = response.Status
		return resp, nil
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		zkLogger.Error(LogTag, "Error while reading the response body: ", err)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return resp, &newZkErr
	}

	x := map[string]dto.TestConnectionResponse{}

	err = json.Unmarshal(bodyBytes, &x)
	if err != nil {
		zkLogger.Error(LogTag, "Error while unmarshalling the response body: ", err)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return resp, &newZkErr
	}

	return x["payload"], nil
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

	return done, common.ToPtr(id), nil
}

func getIntegrationDetails(i integrationsService, integrationId string) ([]dto.Integration, *zkerrors.ZkError) {
	var zkError *zkerrors.ZkError
	integration, err := i.repo.GetAnIntegrationDetails(integrationId)
	if err != nil {
		zkError = common.ToPtr(zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err))
	} else if integration == nil || len(integration) == 0 {
		zkError = common.ToPtr(zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestInvalidClusterAndUrlCombination, err))
	}

	return integration, zkError
}

func getUsernamePassword(integration dto.Integration) (string, string) {
	var auth dto.Auth
	var username, password string
	err := json.Unmarshal(integration.Authentication, &auth)
	if err == nil {
		username = auth.Username
		password = auth.Password
	}

	return username, password
}

func getPrometheusApiResponse(integration dto.Integration) (*http.Response, *zkerrors.ZkError) {
	if common.IsEmpty(integration.ClusterId) || common.IsEmpty(integration.URL) {
		zkLogger.Error(LogTag, "ClusterId or url is empty")
		zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestInvalidClusterAndUrlCombination, nil)
		return nil, &zkError
	}

	url := fmt.Sprintf("http://zk-axon.zk-client.svc.cluster.local:80/v1/c/axon/prom/%s/status", *integration.ID)

	return zkHttp.Create().
		Header("X-PROXY-DESTINATION", url).
		Header("X-CLIENT-ID", integration.ClusterId).
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
