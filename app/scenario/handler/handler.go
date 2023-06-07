package handler

import (
	"github.com/kataras/iris/v12"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	"main/app/cluster/validation"
	"main/app/scenario/service"
	"main/app/scenario/transformer"
	"strconv"
)

type ScenarioHandler interface {
	GetAllScenario(ctx iris.Context)
}

type scenarioHandler struct {
	service service.ScenarioService
}

func NewScenarioHandler(s service.ScenarioService) ScenarioHandler {
	return &scenarioHandler{service: s}
}

func (r scenarioHandler) GetAllScenario(ctx iris.Context) {
	clusterId := ctx.GetHeader("Cluster-Id")
	version := ctx.URLParam("version")
	deleted := ctx.URLParamDefault("deleted", "false")
	limit := ctx.URLParamDefault("limit", "100000")
	offset := ctx.URLParamDefault("offset", "0")
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
