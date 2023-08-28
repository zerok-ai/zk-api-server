package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	model2 "zk-api-server/app/scenario/model"
	"zk-api-server/app/scenario/repository"
	"zk-api-server/app/scenario/transformer"
)

var LogTag = "scenario_service"

type ScenarioService interface {
	GetAllScenarioForOperator(clusterId string, version int64, deleted bool, offset, limit int) (*transformer.ScenarioResponse, *zkerrors.ZkError)
	GetAllScenarioForDashboard(clusterId string, version int64, deleted bool, offset, limit int) (*transformer.ScenarioResponse, *zkerrors.ZkError)
	CreateScenario(clusterId string, request model2.CreateScenarioRequest) *zkerrors.ZkError
	DeleteScenario(clusterId, scenarioId string) *zkerrors.ZkError
	DisableScenario(clusterId, scenarioId string) *zkerrors.ZkError
}

type scenarioService struct {
	repo repository.ScenarioRepo
}

func NewScenarioService(repo repository.ScenarioRepo) ScenarioService {
	return &scenarioService{repo: repo}
}

func (r scenarioService) CreateScenario(clusterId string, request model2.CreateScenarioRequest) *zkerrors.ZkError {
	err := r.repo.CreateNewScenario(clusterId, request)
	if err != nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return &zkError
	}
	return nil
}

func (r scenarioService) GetAllScenarioForOperator(clusterId string, version int64, deleted bool, offset, limit int) (*transformer.ScenarioResponse, *zkerrors.ZkError) {
	scenarioList, deletedScenariosList, disabledScenariosList, zkErr := getAllScenarioData(r.repo, clusterId, version, deleted, offset, limit)
	if zkErr != nil {
		return nil, zkErr
	}

	var deletedScenarioIdList, disabledScenarioIdList []string

	for _, s := range deletedScenariosList {
		deletedScenarioIdList = append(deletedScenarioIdList, s.Scenario.Id)
	}

	for _, s := range disabledScenariosList {
		disabledScenarioIdList = append(disabledScenarioIdList, s.Scenario.Id)
	}

	return transformer.FromScenarioArrayToScenarioResponse(&scenarioList, &deletedScenarioIdList, &disabledScenarioIdList), nil
}

func (r scenarioService) GetAllScenarioForDashboard(clusterId string, version int64, deleted bool, offset, limit int) (*transformer.ScenarioResponse, *zkerrors.ZkError) {
	scenarioList, _, disabledScenariosList, zkErr := getAllScenarioData(r.repo, clusterId, version, deleted, offset, limit)
	if zkErr != nil {
		return nil, zkErr
	}

	scenarioList = append(scenarioList, disabledScenariosList...)
	return transformer.FromScenarioArrayToScenarioResponse(&scenarioList, nil, nil), nil
}

func getAllScenarioData(repo repository.ScenarioRepo, clusterId string, version int64, deleted bool, offset, limit int) ([]transformer.ScenarioModelResponse, []transformer.ScenarioModelResponse, []transformer.ScenarioModelResponse, *zkerrors.ZkError) {
	filter := repository.ScenarioQueryFilter{
		ClusterId: clusterId,
		Deleted:   deleted,
		Version:   version,
		Limit:     limit,
		Offset:    offset,
	}

	scenarioList, err := repo.GetAllScenario(&filter)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			zkLogger.Error(LogTag, "no rows were returned", err)
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorNotFound, err)
			return nil, nil, nil, &zkError
		case err == nil:
			break
		default:
			zkLogger.Error(LogTag, "some db error occurred", err)
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return nil, nil, nil, &zkError
		}
	}

	var activeScenariosList []transformer.ScenarioModelResponse
	var deletedScenariosList []transformer.ScenarioModelResponse
	var disabledScenariosList []transformer.ScenarioModelResponse
	for _, rs := range *scenarioList {
		var d model.Scenario
		err := json.Unmarshal([]byte(rs.ScenarioData), &d)
		if err != nil || d.Workloads == nil {
			zkLogger.Error(LogTag, err)
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return nil, nil, nil, &zkError
		}

		s := transformer.ScenarioModelResponse{
			Scenario:   d,
			CreatedAt:  rs.CreatedAt,
			DisabledAt: rs.DisabledAt,
		}

		if rs.Deleted == true {
			deletedScenariosList = append(deletedScenariosList, s)
		} else if rs.Disabled == true {
			disabledScenariosList = append(disabledScenariosList, s)
		} else {
			activeScenariosList = append(activeScenariosList, s)
		}
	}

	return activeScenariosList, deletedScenariosList, disabledScenariosList, nil
}

func (r scenarioService) DisableScenario(clusterId, scenarioId string) *zkerrors.ZkError {
	_, err := r.repo.DisableScenario(clusterId, scenarioId)
	if err != nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return &zkError
	}
	return nil
}

func (r scenarioService) DeleteScenario(clusterId, scenarioId string) *zkerrors.ZkError {
	_, err := r.repo.DeleteScenario(clusterId, scenarioId)
	if err != nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return &zkError
	}

	return nil
}
