package transformer

import (
	"github.com/zerok-ai/zk-utils-go/rules/model"
)

type RulesResponse struct {
	Rules []model.FilterRule `json:"rules"`
}

func FromFilterRuleArrayToRulesResponse(rArr []model.FilterRule) *RulesResponse {
	var resp RulesResponse
	if rArr != nil && len(rArr) != 0 {
		for _, v := range rArr {
			resp.Rules = append(resp.Rules, v)
		}
	}
	return &resp
}
