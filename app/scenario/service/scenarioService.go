package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"time"
	model2 "zk-api-server/app/scenario/model"
	"zk-api-server/app/scenario/repository"
	"zk-api-server/app/scenario/transformer"
)

var LogTag = "scenario_service"

type ScenarioService interface {
	GetAllScenarioForOperator(clusterId string, version int64, deleted bool, offset, limit int) (*transformer.ScenarioResponse, *zkerrors.ZkError)
	GetAllScenarioForDashboard(clusterId string, version int64, deleted bool, offset, limit int) (*transformer.ScenarioResponse, *zkerrors.ZkError)
	CreateScenario(clusterId string, request model2.CreateScenarioRequest) *zkerrors.ZkError
	DisableScenario(clusterId, scenarioId string, disable bool) *zkerrors.ZkError
	DeleteScenario(clusterId, scenarioId string) *zkerrors.ZkError
}

type scenarioService struct {
	repo repository.ScenarioRepo
}

func NewScenarioService(repo repository.ScenarioRepo) ScenarioService {
	return &scenarioService{repo: repo}
}

func (r scenarioService) CreateScenario(clusterId string, request model2.CreateScenarioRequest) *zkerrors.ZkError {
	for _, s := range request.Workloads {
		if s.Executor != model.ExecutorEbpf && s.Executor != model.ExecutorOTel {
			zkLogger.ErrorF(LogTag, "Executor is not valid, scenario: %v executor: %v", request.ScenarioTitle, s.Executor)
			ZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, errors.New("executor is not valid"))
			return &ZkErr
		}
	}
	err := r.repo.CreateNewScenario(clusterId, request)
	if err != nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return &zkError
	}
	return nil
}

func (r scenarioService) GetAllScenarioForOperator(clusterId string, version int64, deleted bool, offset, limit int) (*transformer.ScenarioResponse, *zkerrors.ZkError) {
	filter := repository.ScenarioQueryFilter{
		ClusterId: clusterId,
		Deleted:   nil,
		Version:   version,
		Limit:     limit,
		Offset:    offset,
	}

	activeAndDisabledScenariosList, deletedScenariosList, zkErr := getAllScenarioData(r.repo, filter)
	if zkErr != nil {
		return nil, zkErr
	}

	var deletedScenarioIdList, disabledScenarioIdList []string

	for _, s := range deletedScenariosList {
		deletedScenarioIdList = append(deletedScenarioIdList, s.Scenario.Id)
	}

	var activeScenariosList []transformer.ScenarioModelResponse

	for _, s := range activeAndDisabledScenariosList {
		if s.DisabledAt == nil {
			activeScenariosList = append(activeScenariosList, s)
		} else {
			disabledScenarioIdList = append(disabledScenarioIdList, s.Scenario.Id)
		}

	}

	return transformer.FromScenarioArrayToScenarioResponse(&activeScenariosList, &deletedScenarioIdList, &disabledScenarioIdList, 0), nil
}

func (r scenarioService) GetAllScenarioForDashboard(clusterId string, version int64, deleted bool, offset, limit int) (*transformer.ScenarioResponse, *zkerrors.ZkError) {
	filter := repository.ScenarioQueryFilter{
		ClusterId: clusterId,
		Deleted:   common.ToPtr(false),
		Version:   version,
		Limit:     limit,
		Offset:    offset,
	}

	scenarioList, _, zkErr := getAllScenarioData(r.repo, filter)
	if zkErr != nil {
		return nil, zkErr
	}

	totalRows, err := r.repo.GetTotalRowsCount(&filter)
	if err != nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorDbError, err)
		return nil, &zkError
	}

	return transformer.FromScenarioArrayToScenarioResponse(&scenarioList, nil, nil, totalRows), nil
}

func getAllScenarioData(repo repository.ScenarioRepo, filter repository.ScenarioQueryFilter) ([]transformer.ScenarioModelResponse, []transformer.ScenarioModelResponse, *zkerrors.ZkError) {
	scenarioList, err := repo.GetAllScenario(&filter)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			zkLogger.Error(LogTag, "no rows were returned", err)
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorNotFound, err)
			return nil, nil, &zkError
		case err == nil:
			break
		default:
			zkLogger.Error(LogTag, "some db error occurred", err)
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return nil, nil, &zkError
		}
	}

	var activeAndDisabledScenariosList []transformer.ScenarioModelResponse
	var deletedScenariosList []transformer.ScenarioModelResponse
	for _, rs := range *scenarioList {
		var d model.Scenario
		err := json.Unmarshal([]byte(rs.ScenarioData), &d)
		if err != nil || d.Workloads == nil {
			zkLogger.Error(LogTag, err)
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return nil, nil, &zkError
		}

		s := transformer.ScenarioModelResponse{
			Scenario:   d,
			CreatedAt:  rs.CreatedAt,
			DisabledAt: rs.DisabledAt,
			UpdatedAt:  rs.UpdatedAt,
		}

		if rs.Deleted == true {
			deletedScenariosList = append(deletedScenariosList, s)
		} else {
			activeAndDisabledScenariosList = append(activeAndDisabledScenariosList, s)
		}
	}

	return activeAndDisabledScenariosList, deletedScenariosList, nil
}

func (r scenarioService) DisableScenario(clusterId, scenarioId string, disable bool) *zkerrors.ZkError {
	var t *int64
	currTime := time.Now().Unix()
	if disable {
		t = common.ToPtr(currTime)
	} else {
		t = nil
	}
	_, err := r.repo.DisableScenario(clusterId, scenarioId, disable, t, currTime)
	if err != nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return &zkError
	}
	return nil
}

func (r scenarioService) DeleteScenario(clusterId, scenarioId string) *zkerrors.ZkError {
	_, err := r.repo.DeleteScenario(clusterId, time.Now().Unix(), scenarioId)
	if err != nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return &zkError
	}

	return nil
}
