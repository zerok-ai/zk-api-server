package handler

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkIntegration "github.com/zerok-ai/zk-utils-go/integration/model"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strings"
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
	GetIntegrationStatus(ctx iris.Context)
}

type integrationsHandler struct {
	service service.IntegrationsService
	cfg     model.ZkApiServerConfig
}

var LogTag = "integrations_handler"

func NewIntegrationsHandler(s service.IntegrationsService, cfg model.ZkApiServerConfig) IntegrationsHandler {
	return &integrationsHandler{service: s, cfg: cfg}
}

func (i integrationsHandler) GetAllIntegrationsOperator(ctx iris.Context) {
	clusterIdx := ctx.Params().Get(utils.ClusterIdxPathParam)
	if zkCommon.IsEmpty(clusterIdx) {
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
	if zkCommon.IsEmpty(clusterIdx) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("ClusterIdx is required")
		return
	}

	zkHttpResponse := getAllIntegrations(i, false, clusterIdx)

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (i integrationsHandler) GetIntegrationStatus(ctx iris.Context) {
	integrationId := ctx.Params().Get(utils.IntegrationIdxPathParam)
	if zkCommon.IsEmpty(integrationId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("IntegrationId is required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[any]
	var zkErr *zkerrors.ZkError
	statusCode, zkErr := i.service.GetAnIntegrationDetails(integrationId)

	if i.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[any](statusCode, nil, nil, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[any](statusCode, nil, nil, zkErr)
	}

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
	clusterIdx := ctx.Params().Get(utils.ClusterIdxPathParam)
	if zkCommon.IsEmpty(clusterIdx) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("ClusterIdx is required")
		return
	}

	var request dto.IntegrationRequest
	var zkHttpResponse zkHttp.ZkHttpResponse[dto.UpsertIntegrationResponse]

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

	request.ClusterId = clusterIdx
	request.Type = zkIntegration.Type(strings.ToUpper(string(request.Type)))
	request.Level = zkIntegration.Level(strings.ToUpper(string(request.Level)))
	err = validation.ValidateIntegrationsUpsertRequest(request)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	done, insertId, zkError := i.service.UpsertIntegration(transformer.FromIntegrationsRequestToIntegrationsDto(request))

	if done {
		zkHttpResponse.Data.IntegrationId = *insertId
		status, zkError := i.service.GetAnIntegrationDetails(*insertId)
		if zkError != nil {
			zkLogger.Error(LogTag, "Error while getting the integration status: ", zkError)
		} else {
			zkHttpResponse.Data.Status = status
		}
	}

	if i.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[dto.UpsertIntegrationResponse](200, zkHttpResponse.Data, done, zkError)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[dto.UpsertIntegrationResponse](200, zkHttpResponse.Data, nil, zkError)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}
