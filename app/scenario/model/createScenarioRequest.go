package model

import (
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"strconv"
)

type GroupByItem struct {
	WorkloadIndex int    `json:"workload_index"`
	Title         string `json:"title"`
	Hash          string `json:"hash"`
}

type CreateScenarioRequest struct {
	ScenarioTitle string            `json:"scenario_title"`
	ScenarioType  string            `json:"scenario_type"`
	Workloads     []model.Workload  `json:"workloads"`
	GroupBy       []GroupByItem     `json:"group_by"`
	RateLimit     []model.RateLimit `json:"rate_limit"`
}

type CreateScenarioResponse struct {
}

type ScenarioState struct {
	Action string `json:"action"`
}

func (cs CreateScenarioRequest) CreateScenarioObj(scenarioId int) model.Scenario {
	var workloadMap = make(map[string]model.Workload)
	currentTime := common.CurrentTime()
	epochTime := currentTime.Unix()
	epochTimeString := strconv.FormatInt(epochTime, 10)
	var workloadIds []string

	for _, workload := range cs.Workloads {
		id := model.WorkLoadUUID(workload)
		workloadMap[id.String()] = workload
		workloadIds = append(workloadIds, id.String())
	}

	defaultFilter := model.Filter{
		Type:        "workload",
		Condition:   "AND",
		WorkloadIds: (*model.WorkloadIds)(&workloadIds),
	}

	// TODO: remove this code once we have a UI to create rate limit
	defaultRateLimit := model.RateLimit{
		BucketMaxSize:    5,
		BucketRefillSize: 5,
		TickDuration:     "1m",
	}

	rateLimitArr := make([]model.RateLimit, 0)

	if len(cs.RateLimit) == 0 {
		rateLimitArr = append(rateLimitArr, defaultRateLimit)
	} else {
		rateLimitArr = append(rateLimitArr, cs.RateLimit...)
	}

	var finalGroupBy []model.GroupBy

	for _, groupBy := range cs.GroupBy {
		newGroupBy := model.GroupBy{
			WorkloadId: workloadIds[groupBy.WorkloadIndex],
			Title:      groupBy.Title,
			Hash:       groupBy.Hash,
		}
		finalGroupBy = append(finalGroupBy, newGroupBy)
	}

	return model.Scenario{
		Version:   epochTimeString,
		Id:        strconv.Itoa(scenarioId),
		Title:     cs.ScenarioTitle,
		Type:      cs.ScenarioType,
		Enabled:   true,
		Workloads: &workloadMap,
		Filter:    defaultFilter,
		RateLimit: rateLimitArr,
		GroupBy:   finalGroupBy,
	}
}

type ScenarioVersionInsertParams struct {
	ScenarioId      int
	ScenarioData    string
	SchemaVersion   string
	ScenarioVersion int64
	CreatedAt       int64
	CreatedBy       string
}

func (si ScenarioVersionInsertParams) GetAllColumns() []any {
	return []any{si.ScenarioId, si.ScenarioData, si.SchemaVersion, si.ScenarioVersion, si.CreatedBy, si.CreatedAt}
}
