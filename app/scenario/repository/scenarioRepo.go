package repository

import (
	"database/sql"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	scenarioResponseModel "main/app/scenario/model"
)

const (
	DefaultClusterId           = "Zk_default_cluster_id_for_all_scenarios"
	GetAllScenarioSqlStatement = `SELECT scenario_data, deleted, disabled FROM scenario s INNER JOIN scenario_version sv USING(scenario_id) WHERE (scenario_version>$1 OR deleted_at>$2 OR disabled_at>$3) AND (cluster_id=$4 OR cluster_id=$5)`
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
	GetAllScenario(filters *ScenarioQueryFilter) (*[]scenarioResponseModel.ScenarioDbResponse, error)
}

type zkPostgresRepo struct {
	dbRepo sqlDB.DatabaseRepo
}

func NewZkPostgresRepo(db sqlDB.DatabaseRepo) ScenarioRepo {
	return &zkPostgresRepo{db}
}

func (zkPostgresRepo zkPostgresRepo) GetAllScenario(filters *ScenarioQueryFilter) (*[]scenarioResponseModel.ScenarioDbResponse, error) {
	params := []any{filters.Version, filters.Version, filters.Version, filters.ClusterId, DefaultClusterId}
	rows, err, closeRow := zkPostgresRepo.dbRepo.GetAll(GetAllScenarioSqlStatement, params)

	return Processor(rows, err, closeRow)
}

func Processor(rows *sql.Rows, sqlErr error, f func()) (*[]scenarioResponseModel.ScenarioDbResponse, error) {
	defer f()

	if sqlErr != nil {
		return nil, sqlErr
	}

	if rows == nil {
		zkLogger.Debug(LogTag, "rows nil", sqlErr)
		return nil, sqlErr
	}

	var scenarioResponse scenarioResponseModel.ScenarioDbResponse
	var scenarioResponseArr []scenarioResponseModel.ScenarioDbResponse
	for rows.Next() {
		err := rows.Scan(&scenarioResponse.ScenarioData, &scenarioResponse.Deleted, &scenarioResponse.Disabled)
		if err != nil {
			zkLogger.Error(LogTag, err)
		}

		scenarioResponseArr = append(scenarioResponseArr, scenarioResponse)
	}

	err := rows.Err()
	if err != nil {
		zkLogger.Error(LogTag, err)
	}

	return &scenarioResponseArr, nil
}
