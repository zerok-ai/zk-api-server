package handler

import (
	"github.com/kataras/iris/v12"
	"main/app/cluster/validation"
	"main/app/ruleengine/service"
	"main/app/ruleengine/transformer"
	"main/app/utils"
	"strconv"
)

type RuleHandler interface {
	GetAllRules(ctx iris.Context)
}

type ruleHandler struct {
	service service.RuleService
}

func NewRuleHandler(s service.RuleService) RuleHandler {
	return &ruleHandler{service: s}
}

func (r ruleHandler) GetAllRules(ctx iris.Context) {
	clusterId := ctx.GetHeader("cluster_id")
	version := ctx.URLParam("version")
	deleted := ctx.URLParamDefault("deleted", "false")
	limit := ctx.URLParamDefault("limit", "100000")
	offset := ctx.URLParamDefault("offset", "0")
	if err := validation.ValidateGetAllRulesApi(clusterId, version, deleted, offset, limit); err != nil {
		z := &utils.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	v, _ := strconv.ParseInt(version, 10, 64)
	d, _ := strconv.ParseBool(deleted)
	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	retVal, err := r.service.GetAllRules(clusterId, v, d, o, l)
	utils.SetResponseInCtxAndReturn[transformer.RulesResponse](ctx, retVal, err)
}
