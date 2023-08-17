package handler

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"strconv"
	"zk-api-server/app/cluster/validation"
	model2 "zk-api-server/app/scenario/model"
	"zk-api-server/app/scenario/service"
	"zk-api-server/app/scenario/transformer"
	"zk-api-server/app/utils"
)

var LogTag = "scenario_handler"

type ScenarioHandler interface {
	GetAllScenario(ctx iris.Context)
	CreateScenario(ctx iris.Context)
}

type scenarioHandler struct {
	service service.ScenarioService
}

func (r scenarioHandler) CreateScenario(ctx iris.Context) {
	defer func() {
		if r := recover(); r != nil {
			zkLogger.Error(LogTag, "Recovered from panic ", r)
			//Send 500 response.
		}
	}()
	clusterId := ctx.Params().Get("clusterIdx")
	zkLogger.Debug(LogTag, "ClusterId is ", clusterId)
	var request model2.CreateScenarioRequest

	// Get the request body as []byte
	body, err := ctx.GetBody()
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Error reading request body")
		return
	}

	// Unmarshal the JSON request body into the struct
	err = json.Unmarshal(body, &request)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Error decoding JSON")
		return
	}
	zkError := r.service.CreateScenario(clusterId, request)
	resp := model2.CreateScenarioResponse{}
	zkHttpResponse := zkHttp.ToZkResponse[model2.CreateScenarioResponse](200, resp, nil, zkError)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)

}

func NewScenarioHandler(s service.ScenarioService) ScenarioHandler {
	return &scenarioHandler{service: s}
}

func (r scenarioHandler) GetAllScenario(ctx iris.Context) {
	clusterId := ctx.GetHeader(utils.ClusterIdHeader)
	version := ctx.URLParam(utils.LastSyncTS)
	deleted := ctx.URLParamDefault(utils.Deleted, "false")
	limit := ctx.URLParamDefault(utils.Limit, "100000")
	offset := ctx.URLParamDefault(utils.Offset, "0")
	if err := validation.ValidateGetAllScenarioApi(clusterId, version, deleted, offset, limit); err != nil {
		zkHttpResponse := zkHttp.ZkHttpResponseBuilder[any]{}.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	v, _ := strconv.ParseInt(version, 10, 64)
	d, _ := strconv.ParseBool(deleted)
	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, zkError := r.service.GetAllScenario(clusterId, v, d, o, l)
	zkHttpResponse := zkHttp.ToZkResponse[transformer.ScenarioResponse](200, *resp, resp, zkError)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)

}
