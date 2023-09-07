package handler

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"zk-api-server/app/integrations/model/dto"
	"zk-api-server/app/integrations/model/transformer"
	"zk-api-server/app/integrations/service"
	"zk-api-server/app/integrations/validation"
	"zk-api-server/app/utils"
	"zk-api-server/internal/model"
)

//var LogTag = "integrations_handler"

type IntegrationsHandler interface {
	GetAllIntegrationsOperator(ctx iris.Context)
	GetAllIntegrationsDashboard(ctx iris.Context)
	UpsertIntegration(ctx iris.Context)
}

type integrationsHandler struct {
	service service.IntegrationsService
	cfg     model.ZkApiServerConfig
}

func NewIntegrationsHandler(s service.IntegrationsService, cfg model.ZkApiServerConfig) IntegrationsHandler {
	return &integrationsHandler{service: s, cfg: cfg}
}

func (i integrationsHandler) GetAllIntegrationsOperator(ctx iris.Context) {
	clusterIdx := ctx.Params().Get(utils.ClusterIdxPathParam)
	if common.IsEmpty(clusterIdx) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("ClusterIdx is required")
		return
	}

	zkHttpResponse := getAllIntegrations(i, true, clusterIdx)

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (i integrationsHandler) GetAllIntegrationsDashboard(ctx iris.Context) {
	clusterIdx := ctx.Params().Get(utils.ClusterIdxPathParam)
	if common.IsEmpty(clusterIdx) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("ClusterIdx is required")
		return
	}

	zkHttpResponse := getAllIntegrations(i, false, clusterIdx)

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func getAllIntegrations(i integrationsHandler, onlyActive bool, clusterId string) zkHttp.ZkHttpResponse[transformer.IntegrationResponse] {

	var zkHttpResponse zkHttp.ZkHttpResponse[transformer.IntegrationResponse]
	var zkErr *zkerrors.ZkError
	var resp transformer.IntegrationResponse

	resp, zkErr = i.service.GetAllIntegrations(clusterId, onlyActive)

	if i.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.IntegrationResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.IntegrationResponse](200, resp, nil, zkErr)
	}

	return zkHttpResponse
}

func (i integrationsHandler) UpsertIntegration(ctx iris.Context) {
	var request dto.IntegrationRequest
	var zkHttpResponse zkHttp.ZkHttpResponse[transformer.IntegrationResponse]

	body, err := ctx.GetBody()
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Error reading request body")
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Error decoding JSON")
		return
	}

	err = validation.ValidateIntegrationsUpsertRequest(request)
	if err != nil {
		return
	}

	done, zkError := i.service.UpsertIntegration(transformer.FromIntegrationsRequestToIntegrationsDto(request))
	if i.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.IntegrationResponse](200, transformer.IntegrationResponse{}, done, zkError)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.IntegrationResponse](200, transformer.IntegrationResponse{}, nil, zkError)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}
