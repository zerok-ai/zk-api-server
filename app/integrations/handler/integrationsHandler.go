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
	TestIntegrationConnectionStatus(ctx iris.Context)
	TestUnSyncedIntegrationConnection(ctx iris.Context)
	GetIntegrationsById(ctx iris.Context)
	DeleteIntegrationById(ctx iris.Context)
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

func (i integrationsHandler) GetIntegrationsById(ctx iris.Context) {
	integrationId := ctx.Params().Get(utils.IntegrationIdxPathParam)
	if zkCommon.IsEmpty(integrationId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("IntegrationId is required")
		return
	}

	clusterIdx := ctx.Params().Get(utils.ClusterIdxPathParam)
	if zkCommon.IsEmpty(clusterIdx) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("ClusterIdx is required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[zkIntegration.IntegrationResponseObj]
	var zkErr *zkerrors.ZkError
	var resp zkIntegration.IntegrationResponseObj

	resp, zkErr = i.service.GetIntegrationById(clusterIdx, integrationId)

	if i.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[zkIntegration.IntegrationResponseObj](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[zkIntegration.IntegrationResponseObj](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (i integrationsHandler) DeleteIntegrationById(ctx iris.Context) {
	integrationId := ctx.Params().Get(utils.IntegrationIdxPathParam)
	if zkCommon.IsEmpty(integrationId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("IntegrationId is required")
		return
	}

	clusterIdx := ctx.Params().Get(utils.ClusterIdxPathParam)
	if zkCommon.IsEmpty(clusterIdx) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("ClusterIdx is required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[any]
	var zkErr *zkerrors.ZkError
	var resp zkIntegration.IntegrationResponseObj

	zkErr = i.service.DeleteIntegrationById(clusterIdx, integrationId)

	if i.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[any](204, nil, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[any](204, nil, nil, zkErr)
	}

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

func (i integrationsHandler) TestIntegrationConnectionStatus(ctx iris.Context) {
	integrationId := ctx.Params().Get(utils.IntegrationIdxPathParam)
	if zkCommon.IsEmpty(integrationId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("IntegrationId is required")
		return
	}

	clusterId := ctx.Params().Get(utils.ClusterIdxPathParam)
	if zkCommon.IsEmpty(clusterId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("ClusterIdx is required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[any]
	var zkErr *zkerrors.ZkError
	resp, zkErr := i.service.TestIntegrationConnection(integrationId, clusterId)

	if i.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[any](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[any](200, resp, resp, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (i integrationsHandler) TestUnSyncedIntegrationConnection(ctx iris.Context) {
	clusterIdx := ctx.Params().Get(utils.ClusterIdxPathParam)
	if zkCommon.IsEmpty(clusterIdx) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("ClusterIdx is required")
		return
	}

	var request dto.IntegrationRequest
	var zkHttpResponse zkHttp.ZkHttpResponse[dto.TestConnectionResponse]

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

	integration := transformer.FromIntegrationsRequestToIntegrationsDto(request)

	resp, zkError := i.service.TestUnSyncedIntegrationConnection(integration)
	if zkError != nil {
		zkLogger.Error(LogTag, "Error while getting the integration status: ", zkError)
	} else {
		zkHttpResponse.Data = resp
	}

	if i.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[dto.TestConnectionResponse](200, zkHttpResponse.Data, zkHttpResponse.Data, zkError)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[dto.TestConnectionResponse](200, zkHttpResponse.Data, zkHttpResponse.Data, zkError)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func getAllIntegrations(i integrationsHandler, forOperator bool, clusterId string) zkHttp.ZkHttpResponse[transformer.IntegrationResponse] {

	var zkHttpResponse zkHttp.ZkHttpResponse[transformer.IntegrationResponse]
	var zkErr *zkerrors.ZkError
	var resp transformer.IntegrationResponse

	resp, zkErr = i.service.GetAllIntegrations(clusterId, forOperator)

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

	integration := transformer.FromIntegrationsRequestToIntegrationsDto(request)
	resp, zkErrorTestConnection := i.service.TestUnSyncedIntegrationConnection(integration)
	integration.MetricServer = resp.IntegrationStatus.HasMetricServer
	done, insertId, zkError := i.service.UpsertIntegration(integration)
	if done {
		zkHttpResponse.Data.IntegrationId = *insertId
		if zkErrorTestConnection != nil {
			zkHttpResponse.Data.IntegrationStatus.ConnectionStatus = utils.StatusError
			zkLogger.Error(LogTag, "Error while getting the integration status: ", zkError)
		} else {
			zkHttpResponse.Data.IntegrationStatus = resp.IntegrationStatus
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
