package repository

import (
	"database/sql"
	"encoding/json"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	zkUtilsPostgres "github.com/zerok-ai/zk-utils-go/storage/sqlDB/postgres"
	scenarioResponseModel "main/app/scenario/model"
)

var LogTag = "scenario_repo"

type ScenarioQueryFilter struct {
	ClusterId string
	Deleted   bool
	Version   int64
	Limit     int
	Offset    int
}

type ScenarioRepo interface {
	GetAllScenario(filters *ScenarioQueryFilter) (*[]model.Scenario, *[]string, error)
}

type zkPostgresRepo struct {
}

func NewZkPostgresRepo() ScenarioRepo {
	return &zkPostgresRepo{}
}

func (zkPostgresService zkPostgresRepo) GetAllScenario(filters *ScenarioQueryFilter) (*[]model.Scenario, *[]string, error) {
	query := GetAllScenarioSqlStatement
	dbRepo := zkUtilsPostgres.NewZkPostgresRepo()
	db, err := dbRepo.GetDBInstance()
	if err != nil {
		zkLogger.Error(LogTag, "unable to get db instance", err)
	}

	params := []any{filters.ClusterId, filters.Version, filters.Limit, filters.Offset}
	rows, err, closeRow := dbRepo.GetAll(db, query, params)

	return Processor(rows, err, closeRow)
}

func Processor(rows *sql.Rows, sqlErr error, f func()) (*[]model.Scenario, *[]string, error) {
	defer f()

	if sqlErr != nil {
		return nil, nil, sqlErr
	}

	if rows == nil {
		zkLogger.Debug(LogTag, "rows nil", sqlErr)
		return nil, nil, sqlErr
	}

	var scenarioResponse scenarioResponseModel.ScenarioDbResponse
	var scenarioResponseArr []scenarioResponseModel.ScenarioDbResponse
	for rows.Next() {
		err := rows.Scan(&scenarioResponse.Scenario, &scenarioResponse.Deleted)
		if err != nil {
			zkLogger.Error(LogTag, err)
		}

		scenarioResponseArr = append(scenarioResponseArr, scenarioResponse)
	}

	// Check for any errors occurred during iteration
	err := rows.Err()
	if err != nil {
		zkLogger.Error(LogTag, err)
	}

	var scenarios []model.Scenario
	var deletedScenarioIdList []string
	for _, rs := range scenarioResponseArr {
		var d model.Scenario
		err := json.Unmarshal([]byte(rs.Scenario), &d)
		if err != nil || d.Workloads == nil {
			zkLogger.Error(LogTag, err)
			return nil, nil, err
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
