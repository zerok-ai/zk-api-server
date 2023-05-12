package handler

import (
	"github.com/kataras/iris/v12"
	"main/app/ruleengine/service"
	"main/app/ruleengine/transformer"
	"main/app/utils"
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
	retVal, err := r.service.GetAllRules()
	utils.SetResponseInCtxAndReturn[transformer.RulesResponse](ctx, retVal, err)
}
