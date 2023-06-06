package repository

import (
	"database/sql"
	"encoding/json"
	"github.com/zerok-ai/zk-utils-go/rules/model"
	"log"
	scenarioResponseModel "main/app/scenario/model"
	"main/app/utils"
	zkLogger "main/app/utils/logs"
	zkPostgres "main/app/utils/postgres"
	"main/app/utils/zkerrors"
)

var LOG_TAG = "zkpostgres_db_repo"

type ScenarioQueryFilter struct {
	ClusterId string
	Deleted   bool
	Version   int64
	Limit     int
	Offset    int
}

type ScenarioRepo interface {
	GetAllScenario(filters *ScenarioQueryFilter) (*[]model.Scenario, *[]string, *zkerrors.ZkError)
}

type zkPostgresRepo struct {
}

func NewZkPostgresRepo() ScenarioRepo {
	return &zkPostgresRepo{}
}

func (zkPostgresService zkPostgresRepo) GetAllScenario(filters *ScenarioQueryFilter) (*[]model.Scenario, *[]string, *zkerrors.ZkError) {
	query := GetAllScenarioSqlStatement
	zkPostgresRepo := zkPostgres.NewZkPostgresRepo[model.Scenario]()

	params := []any{filters.ClusterId, filters.Version, filters.Limit, filters.Offset}
	return zkPostgresRepo.GetAll(query, params, Processor)
}

func Processor(rows *sql.Rows, sqlErr error) (*[]model.Scenario, *[]string, *zkerrors.ZkError) {
	defer rows.Close()

	switch sqlErr {
	case sql.ErrNoRows:
		zkLogger.Debug(LOG_TAG, "no rows were returned", sqlErr)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_NOT_FOUND, sqlErr)
		return nil, nil, &zkError
	case nil:
		break
	default:
		zkLogger.Debug(LOG_TAG, "unable to scan rows", sqlErr)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, sqlErr)
		return nil, nil, &zkError
	}

	if rows == nil {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
		return nil, nil, &zkErr
	}

	var scenarioResponse scenarioResponseModel.ScenarioDbResponse
	var scenarioResponseArr []scenarioResponseModel.ScenarioDbResponse
	for rows.Next() {
		err := rows.Scan(&scenarioResponse.Scenario, &scenarioResponse.Deleted)
		if err != nil {
			log.Fatal(err)
		}

		scenarioResponseArr = append(scenarioResponseArr, scenarioResponse)
	}

	// Check for any errors occurred during iteration
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	var scenarios []model.Scenario
	var deletedScenarioIdList []string
	for _, rs := range scenarioResponseArr {
		var d model.Scenario
		err := json.Unmarshal([]byte(rs.Scenario), &d)
		if err != nil || d.Workloads == nil {
			log.Println(err)
			return nil, nil, utils.ToPtr[zkerrors.ZkError](zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil))
		}

		if rs.Deleted == false {
			scenarios = append(scenarios, d)
			//for oldId, v := range d.Workloads {
			//	id := model.WorkLoadUUID(v)
			//	delete(d.Workloads, oldId)
			//	d.Workloads[id.String()] = v
			//}
		} else {
			deletedScenarioIdList = append(deletedScenarioIdList, d.ScenarioId)
		}
	}

	return &scenarios, &deletedScenarioIdList, nil
}

const GetAllScenarioSqlStatement = `SELECT scenario, deleted FROM Scenario WHERE (cluster_id=$1 OR cluster_id IS NULL) AND version>$2 LIMIT $3 OFFSET $4`
