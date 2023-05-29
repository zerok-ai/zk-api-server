package service

import (
	"log"
	"main/app/scenario/repository"
	"main/app/scenario/transformer"
	"main/app/utils/zkerrors"
)

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

	scenarioList, deletedScenarioIds, zkErr := r.repo.GetAllScenario(&filter)
	if zkErr != nil {
		log.Println(zkErr)
		return nil, zkErr
	}

	return transformer.FromScenarioArrayToScenarioResponse(scenarioList, deletedScenarioIds), nil
}
