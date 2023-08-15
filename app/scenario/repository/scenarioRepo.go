package repository

import (
	"database/sql"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	scenarioResponseModel "zk-api-server/app/scenario/model"
	"zk-api-server/app/utils/errors"
)

const (
	DefaultClusterId             = "Zk_default_cluster_id_for_all_scenarios"
	GetAllScenarioSqlStatement   = `SELECT scenario_data, deleted, disabled FROM scenario s INNER JOIN scenario_version sv USING(scenario_id) WHERE (scenario_version>$1 OR deleted_at>$2 OR disabled_at>$3) AND (cluster_id=$4 OR cluster_id=$5)`
	InsertScenarioTableStatement = "INSERT INTO scenario (cluster_id, scenario_title, scenario_type) VALUES ($1, $2, $3) RETURNING scenario_id"
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
	CreateNewScenario(clusterId string, request scenarioResponseModel.CreateScenarioRequest) error
}

type zkPostgresRepo struct {
	dbRepo sqlDB.DatabaseRepo
}

func NewZkPostgresRepo(db sqlDB.DatabaseRepo) ScenarioRepo {
	return &zkPostgresRepo{db}
}

func (zkPostgresRepo zkPostgresRepo) CreateNewScenario(clusterId string, request scenarioResponseModel.CreateScenarioRequest) error {
	tx, err := zkPostgresRepo.dbRepo.CreateTransaction()

	if err != nil {
		zkLogger.Error(LogTag, "Error while creating a db transaction in createNewScenario ", err)
		return errors.ErrInternalServerError
	}

	params := []any{clusterId, request.ScenarioTitle, request.ScenarioType}

	var scenarioId string

	//TODO: Discuss with vaibhav about this.
	err = zkPostgresRepo.dbRepo.Get(InsertScenarioTableStatement, params, []any{&scenarioId})

	if err != nil {
		zkLogger.Error(LogTag, "Error while executing the insert query ", err)
		return err
	}

	zkLogger.Debug(LogTag, "New scenarioId is ", scenarioId)

	err = tx.Commit()
	if err != nil {
		zkLogger.Error(LogTag, "Error while committing a db transaction in createNewScenario ", err)
		return err
	}
	zkLogger.Debug(LogTag, "Reached the end of the handler method.")
	return nil
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
