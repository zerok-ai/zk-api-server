package handler

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
	"zk-api-server/app/cluster/validation"
	scenarioModel "zk-api-server/app/scenario/model"
	"zk-api-server/app/scenario/service"
	"zk-api-server/app/scenario/transformer"
	"zk-api-server/app/utils"
	"zk-api-server/internal/model"
)

var LogTag = "scenario_handler"

type ScenarioHandler interface {
	GetAllScenarioOperator(ctx iris.Context)
	GetAllScenarioDashboard(ctx iris.Context)
	GetScenarioByIdDashboard(ctx iris.Context)
	CreateScenario(ctx iris.Context)
	UpdateScenarioState(ctx iris.Context)
	DeleteScenario(ctx iris.Context)
}

type scenarioHandler struct {
	service service.ScenarioService
	cfg     model.ZkApiServerConfig
}

func NewScenarioHandler(s service.ScenarioService, cfg model.ZkApiServerConfig) ScenarioHandler {
	return &scenarioHandler{service: s, cfg: cfg}
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
	var request scenarioModel.CreateScenarioRequest

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

	var zkHttpResponse zkHttp.ZkHttpResponse[scenarioModel.CreateScenarioResponse]
	var zkErr *zkerrors.ZkError
	var resp scenarioModel.CreateScenarioResponse

	zkErr = r.service.CreateScenario(clusterId, request)

	if r.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[scenarioModel.CreateScenarioResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[scenarioModel.CreateScenarioResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (r scenarioHandler) GetAllScenarioOperator(ctx iris.Context) {
	defer func() {
		if r := recover(); r != nil {
			zkLogger.Error(LogTag, "Recovered from panic ", r)
			//Send 500 response.
		}
	}()
	getAllScenarioHelper(r.service, ctx, false, r.cfg.Http.Debug)

}

func (r scenarioHandler) GetAllScenarioDashboard(ctx iris.Context) {
	// TODO: ask rajeev why we have put this block
	//defer func() {
	//	if r := recover(); r != nil {
	//		zkLogger.Error(LogTag, "Recovered from panic ", r)
	//		//Send 500 response.
	//	}
	//}()
	getAllScenarioHelper(r.service, ctx, true, r.cfg.Http.Debug)
}

func (r scenarioHandler) GetScenarioByIdDashboard(ctx iris.Context) {
	clusterId := ctx.Params().Get(utils.ClusterIdxPathParam)
	scenarioId := ctx.Params().Get(utils.ScenarioIdxPathParam)

	var zkHttpResponse zkHttp.ZkHttpResponse[transformer.ScenarioModelResponse]
	var zkErr *zkerrors.ZkError
	var resp transformer.ScenarioModelResponse

	if zkErr = validation.ValidateGetScenarioByIdApi(clusterId, scenarioId); zkErr != nil {
		zkLogger.Error(LogTag, "Error validating get scenario by Id api ", zkErr)
	} else {
		resp, zkErr = r.service.GetScenarioByIdForDashboard(clusterId, scenarioId)
	}

	if r.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ScenarioModelResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ScenarioModelResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (r scenarioHandler) UpdateScenarioState(ctx iris.Context) {
	defer func() {
		if r := recover(); r != nil {
			zkLogger.Error(LogTag, "Recovered from panic ", r)
			//Send 500 response.
		}
	}()
	clusterId := ctx.Params().Get(utils.ClusterIdxPathParam)
	scenarioId := ctx.Params().Get(utils.ScenarioIdxPathParam)

	var request scenarioModel.ScenarioState

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

	var zkErr *zkerrors.ZkError

	disable := true
	if request.Action == utils.Enable {
		disable = false
	}

	if zkErr = validation.ValidateDisableScenarioApi(clusterId, scenarioId, request); zkErr != nil {
		zkLogger.Error(LogTag, "Error validating disable scenario api ", zkErr)
	} else {
		zkErr = r.service.DisableScenario(clusterId, scenarioId, disable)
	}

	if zkErr != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	ctx.StatusCode(iris.StatusOK)
	return
}

func (r scenarioHandler) DeleteScenario(ctx iris.Context) {
	defer func() {
		if r := recover(); r != nil {
			zkLogger.Error(LogTag, "Recovered from panic ", r)
			//Send 500 response.
		}
	}()
	clusterId := ctx.Params().Get(utils.ClusterIdxPathParam)
	scenarioId := ctx.Params().Get(utils.ScenarioIdxPathParam)

	if err := validation.ValidateDeleteScenarioApi(clusterId, scenarioId); err != nil {
		zkLogger.Error(LogTag, "Error validating delete scenario api ", err)
		zkHttpResponse := zkHttp.ZkHttpResponseBuilder[any]{}.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	err := r.service.DeleteScenario(clusterId, scenarioId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	ctx.StatusCode(iris.StatusOK)
	return
}

func getAllScenarioHelper(service service.ScenarioService, ctx iris.Context, dashboardCall, debug bool) {
	clusterId := ctx.Params().Get(utils.ClusterIdxPathParam)
	version := ctx.URLParam(utils.LastSyncTS)
	deleted := ctx.URLParamDefault(utils.Deleted, "false")
	limit := ctx.URLParamDefault(utils.Limit, "10000")
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

	var resp transformer.ScenarioResponse
	var zkErr *zkerrors.ZkError
	var zkHttpResponse zkHttp.ZkHttpResponse[transformer.ScenarioResponse]

	if dashboardCall {
		resp, zkErr = service.GetAllScenarioForDashboard(clusterId, v, d, o, l)
	} else {
		resp, zkErr = service.GetAllScenarioForOperator(clusterId, v, d, o, l)
	}

	if debug {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ScenarioResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ScenarioResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}
