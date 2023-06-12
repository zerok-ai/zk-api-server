package transformer

import (
	"github.com/zerok-ai/zk-utils-go/scenario/model"
)

type ScenarioResponse struct {
	Scenarios       []model.Scenario `json:"scenarios"`
	DeletedFilterId []string         `json:"deleted_scenario_id"`
}

func FromScenarioArrayToScenarioResponse(sArr *[]model.Scenario, deletedIdArr *[]string) *ScenarioResponse {
	var resp ScenarioResponse
	if sArr != nil && len(*sArr) != 0 {
		for _, v := range *sArr {
			resp.Scenarios = append(resp.Scenarios, v)
		}
	} else {
		resp.Scenarios = make([]model.Scenario, 0)
	}

	if deletedIdArr != nil && len(*deletedIdArr) != 0 {
		resp.DeletedFilterId = *deletedIdArr
	} else {
		resp.DeletedFilterId = make([]string, 0)

	}
	return &resp
}
