package transformer

import (
	"github.com/zerok-ai/zk-utils-go/scenario/model"
)

type ScenarioResponse struct {
	Scenarios        []ScenarioModelResponse `json:"scenarios"`
	DeletedFilterId  []string                `json:"deleted_scenario_id,omitempty"`
	DisabledFilterId []string                `json:"disabled_filter_id,omitempty"`
}

type ScenarioModelResponse struct {
	Scenario   model.Scenario `json:"scenario"`
	CreatedAt  int64          `json:"created_at"`
	DisabledAt *int64         `json:"disabled_at"`
}

func FromScenarioArrayToScenarioResponse(sArr *[]ScenarioModelResponse, deletedIdArr *[]string, disabledIdArr *[]string) *ScenarioResponse {
	var resp ScenarioResponse
	if sArr != nil && len(*sArr) != 0 {
		resp.Scenarios = append(resp.Scenarios, *sArr...)
	} else {
		resp.Scenarios = make([]ScenarioModelResponse, 0)
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
