package service

import (
	"database/sql"
	"encoding/json"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"zk-api-server/app/scenario/repository"
	"zk-api-server/app/scenario/transformer"
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

	scenarioList, err := r.repo.GetAllScenario(&filter)

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

	var scenarios []model.Scenario
	var deletedScenarioIdList []string
	var disabledScenarioIdList []string
	for _, rs := range *scenarioList {
		var d model.Scenario
		err := json.Unmarshal([]byte(rs.ScenarioData), &d)
		if err != nil || d.Workloads == nil {
			zkLogger.Error(LogTag, err)
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
			return nil, &zkError
		}

		if rs.Deleted == true {
			deletedScenarioIdList = append(deletedScenarioIdList, d.Id)
		} else if rs.Disabled == true {
			disabledScenarioIdList = append(disabledScenarioIdList, d.Id)
		} else {
			//workLoadIds := make(model.WorkloadIds, 0)
			//for oldId, v := range *d.Workloads {
			//	id := model.WorkLoadUUID(v)
			//	if oldId != id.String() {
			//		delete(*d.Workloads, oldId)
			//		(*d.Workloads)[id.String()] = v
			//	}
			//	workLoadIds = append(workLoadIds, id.String())
			//}
			//d.Filter.WorkloadIds = &workLoadIds
			scenarios = append(scenarios, d)
		}
	}

	return transformer.FromScenarioArrayToScenarioResponse(&scenarios, &deletedScenarioIdList, &disabledScenarioIdList), nil
}
