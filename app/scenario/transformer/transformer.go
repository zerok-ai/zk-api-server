package transformer

import (
	"github.com/zerok-ai/zk-utils-go/scenario/model"
)

type ScenarioResponse struct {
	Scenarios        []model.Scenario `json:"scenarios"`
	DeletedFilterId  []string         `json:"deleted_scenario_id"`
	DisabledFilterId []string         `json:"disabled_filter_id"`
}

func FromScenarioArrayToScenarioResponse(sArr *[]model.Scenario, deletedIdArr *[]string, disabledIdArr *[]string) *ScenarioResponse {
	var resp ScenarioResponse
	if sArr != nil && len(*sArr) != 0 {
		resp.Scenarios = append(resp.Scenarios, *sArr...)
	} else {
		resp.Scenarios = make([]model.Scenario, 0)
	}

	if deletedIdArr != nil && len(*deletedIdArr) != 0 {
		resp.DeletedFilterId = *deletedIdArr
	} else {
		resp.DeletedFilterId = make([]string, 0)
	}

	if disabledIdArr != nil && len(*disabledIdArr) != 0 {
		resp.DisabledFilterId = *disabledIdArr
	} else {
		resp.DisabledFilterId = make([]string, 0)
	}

	return &resp
}
