package transformer

import (
	"main/app/ruleengine/model"
)

type RulesResponse struct {
	Rules []model.NewRuleSchema `json:"rules"`
}

func FromFilterRuleArrayToRulesResponse(rArr []model.NewRuleSchema) *RulesResponse {
	var resp RulesResponse
	if rArr != nil && len(rArr) != 0 {
		for _, v := range rArr {
			resp.Rules = append(resp.Rules, v)
		}
	}
	return &resp
}
