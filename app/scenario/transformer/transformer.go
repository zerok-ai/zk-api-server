package transformer

import (
	"github.com/zerok-ai/zk-utils-go/scenario/model"
)

type ScenarioResponse struct {
	Scenarios          []ScenarioModelResponse `json:"scenarios"`
	DeletedScenarioId  []string                `json:"deleted_scenario_id,omitempty"`
	DisabledScenarioId []string                `json:"disabled_scenario_id,omitempty"`
	TotalRows          int                     `json:"total_rows,omitempty"`
}

type ScenarioModelResponse struct {
	Scenario   model.Scenario `json:"scenario"`
	CreatedAt  int64          `json:"created_at"`
	DisabledAt *int64         `json:"disabled_at,omitempty"`
	UpdatedAt  int64          `json:"updated_at"`
}

func FromScenarioArrayToScenarioResponse(sArr *[]ScenarioModelResponse, deletedIdArr *[]string, disabledIdArr *[]string, totalRows int) *ScenarioResponse {
	var resp ScenarioResponse
	if sArr != nil && len(*sArr) != 0 {
		resp.Scenarios = append(resp.Scenarios, *sArr...)
	} else {
		resp.Scenarios = make([]ScenarioModelResponse, 0)
	}

	if deletedIdArr != nil && len(*deletedIdArr) != 0 {
		resp.DeletedScenarioId = *deletedIdArr
	} else {
		resp.DeletedScenarioId = make([]string, 0)
	}

	if disabledIdArr != nil && len(*disabledIdArr) != 0 {
		resp.DisabledScenarioId = *disabledIdArr
	} else {
		resp.DisabledScenarioId = make([]string, 0)
	}

	resp.TotalRows = totalRows

	return &resp
}
