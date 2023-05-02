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

// GetAllRules Returns all the rules for RuleEngine processing godoc
//
//	@Summary		Get all rules
//	@Description	Returns all the rules for RuleEngine processing
//	@Tags			rule engine
//	@Produce		json
//	@Success		200 {object} utils.ZkHttpResponse[model.FilterRule]
//	@Router			/o/cluster/rules [get]
func (r ruleHandler) GetAllRules(ctx iris.Context) {
	retVal, err := r.service.GetAllRules()
	utils.SetResponseInCtxAndReturn[transformer.RulesResponse](ctx, retVal, err)
}
