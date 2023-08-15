package repository

import (
	"database/sql"
	"encoding/json"
	"github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	scenarioResponseModel "zk-api-server/app/scenario/model"
	"zk-api-server/app/utils/errors"
)

const (
	DefaultClusterId                    = "Zk_default_cluster_id_for_all_scenarios"
	GetAllScenarioSqlStatement          = `SELECT scenario_data, deleted, disabled FROM scenario s INNER JOIN scenario_version sv USING(scenario_id) WHERE (scenario_version>$1 OR deleted_at>$2 OR disabled_at>$3) AND (cluster_id=$4 OR cluster_id=$5)`
	InsertScenarioTableStatement        = "INSERT INTO scenario (cluster_id, scenario_title, scenario_type) VALUES ($1, $2, $3) RETURNING scenario_id"
	InsertScenarioVersionTableStatement = "INSERT INTO scenario_version (scenario_id, scenario_data, schema_version, scenario_version, created_by, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
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

	var scenarioId int

	//TODO: Discuss with vaibhav about this. Should scenario id state from 1000?
	err = zkPostgresRepo.dbRepo.Get(InsertScenarioTableStatement, params, []any{&scenarioId})

	if err != nil {
		zkLogger.Error(LogTag, "Error while executing the insert query ", err)
		return err
	}

	zkLogger.Debug(LogTag, "New scenarioId is ", scenarioId)

	scenarioObj := request.CreateScenarioObj(scenarioId)
	scenarioData, err := json.Marshal(scenarioObj)

	if err != nil {
		zkLogger.Error(LogTag, "Error while serializing scenario data ", err)
		return err
	}

	currentTime := common.CurrentTime()
	epochTime := currentTime.Unix()

	scenarioVersionParams := scenarioResponseModel.ScenarioVersionInsertParams{
		ScenarioId:      scenarioId,
		ScenarioVersion: epochTime,
		ScenarioData:    string(scenarioData),
		SchemaVersion:   "v1",
		CreatedAt:       epochTime,
		CreatedBy:       "dashboard",
	}

	stmt, err := common.GetStmtRawQuery(tx, InsertScenarioVersionTableStatement)

	if err != nil {
		zkLogger.Error(LogTag, "Error while creating the scenario version insert stmt ", err)
		return err
	}

	_, err = zkPostgresRepo.dbRepo.Insert(stmt, scenarioVersionParams)

	if err != nil {
		zkLogger.Error(LogTag, "Error while inserting into the scenario version table. ", err)
		return err
	}

	//TODO: Discuss with vaibhav that db connection is getting closed after insert.
	done, err2 := common.CommitTransaction(tx, LogTag)
	if err2 != nil {
		zkLogger.Error(LogTag, "Error while committing a db transaction in createNewScenario ", err2.Error)
		return err
	}

	if !done {
		zkLogger.Error(LogTag, "Transaction commit failed. ")
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
