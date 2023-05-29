package transformer

import (
	"github.com/zerok-ai/zk-utils-go/rules/model"
)

type RulesResponse struct {
	Rules           []model.Scenario `json:"rules"`
	DeletedFilterId []string         `json:"deleted_filter_id"`
}

func FromFilterRuleArrayToRulesResponse(rArr *[]model.Scenario, deletedIdArr *[]string) *RulesResponse {
	var resp RulesResponse
	if rArr != nil && len(*rArr) != 0 {
		for _, v := range *rArr {
			resp.Rules = append(resp.Rules, v)
		}
	} else {
		resp.Rules = make([]model.Scenario, 0)
	}

	if deletedIdArr != nil && len(*deletedIdArr) != 0 {
		resp.DeletedFilterId = *deletedIdArr
	} else {
		resp.DeletedFilterId = make([]string, 0)

	}
	return &resp
}
