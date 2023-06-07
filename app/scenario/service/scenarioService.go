package service

import (
	"database/sql"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"main/app/scenario/repository"
	"main/app/scenario/transformer"
)

var LogTag = "scenario_service"

type ScenarioService interface {
	GetAllScenario(clusterId string, version int64, deleted bool, offset, limit int) (*transformer.ScenarioResponse, *zkerrors.ZkError)
}

type scenarioService struct {
	repo repository.ScenarioRepo
}

func NewScenarioService(repo repository.ScenarioRepo) ScenarioService {
	return &scenarioService{repo: repo}
}

func (r scenarioService) GetAllScenario(clusterId string, version int64, deleted bool, offset, limit int) (*transformer.ScenarioResponse, *zkerrors.ZkError) {
	filter := repository.ScenarioQueryFilter{
		ClusterId: clusterId,
		Deleted:   deleted,
		Version:   version,
		Limit:     limit,
		Offset:    offset,
	}

	scenarioList, deletedScenarioIds, err := r.repo.GetAllScenario(&filter)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			zkLogger.Error(LogTag, "no rows were returned", err)
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorNotFound, err)
			return nil, &zkError
		case nil:
			break
		default:
			zkLogger.Error(LogTag, "some db error occurred", err)
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return nil, &zkError
		}
	}

	return transformer.FromScenarioArrayToScenarioResponse(scenarioList, deletedScenarioIds), nil
}
