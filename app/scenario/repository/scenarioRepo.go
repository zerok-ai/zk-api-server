package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	"time"
	scenarioResponseModel "zk-api-server/app/scenario/model"
	"zk-api-server/utils"
)

const (
	DefaultClusterId = "Zk_default_cluster_id_for_all_scenarios"

	GetAllScenarioSqlStatement             = "SELECT scenario_data, deleted, disabled, created_at, disabled_at, updated_at FROM zk_scenario s INNER JOIN zk_scenario_version sv USING(scenario_id) WHERE updated_at>$1 AND (cluster_id=$2 OR cluster_id=$3) order by created_at desc LIMIT $4 OFFSET $5"
	GetScenarioByIdSqlStatement            = "SELECT scenario_data, deleted, disabled, created_at, disabled_at, updated_at FROM zk_scenario s INNER JOIN zk_scenario_version sv USING(scenario_id) WHERE scenario_id>$1 AND (cluster_id=$2 OR cluster_id=$3)"
	GetAllScenarioForDashboardSqlStatement = "SELECT scenario_data, deleted, disabled, created_at, disabled_at, updated_at FROM zk_scenario s INNER JOIN zk_scenario_version sv USING(scenario_id) WHERE updated_at>$1 AND deleted = $2 AND (cluster_id=$3 OR cluster_id=$4) order by created_at desc LIMIT $5 OFFSET $6"
	InsertScenarioTableStatement           = "INSERT INTO zk_scenario (cluster_id, scenario_title, scenario_type, updated_at) VALUES ($1, $2, $3, $4) RETURNING scenario_id"
	InsertScenarioVersionTableStatement    = "INSERT INTO zk_scenario_version (scenario_id, scenario_data, schema_version, scenario_version, created_by, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	DisableScenarioStatement               = "UPDATE zk_scenario set disabled=$1, disabled_at=$2, updated_at=$3 where cluster_id= $4 AND scenario_id=$5"
	DeleteScenarioStatement                = "UPDATE zk_scenario set deleted=TRUE, deleted_at=$1, updated_at=$2 where cluster_id= $3 AND scenario_id=$4"
	GetTotalRowsCountStatement             = "SELECT COUNT(*) as count FROM zk_scenario s INNER JOIN zk_scenario_version sv USING(scenario_id) WHERE deleted=false AND updated_at>$1 AND (cluster_id=$2 OR cluster_id=$3)"
)

var LogTag = "scenario_repo"

type ScenarioQueryFilter struct {
	ClusterId string
	Deleted   *bool
	Version   int64
	Limit     int
	Offset    int
}

type ScenarioRepo interface {
	GetAllScenario(filters *ScenarioQueryFilter) (*[]scenarioResponseModel.ScenarioDbResponse, error)
	GetScenarioById(clusterId, scenarioId string) (*[]scenarioResponseModel.ScenarioDbResponse, error)
	CreateNewScenario(clusterId string, request scenarioResponseModel.CreateScenarioRequest) error
	DisableScenario(clusterId, scenarioId string, disable bool, disabledAtTime *int64, currentTime int64) (int, error)
	DeleteScenario(clusterId string, currentTime int64, scenarioId string) (int, error)
	GetTotalRowsCount(filters *ScenarioQueryFilter) (int, error)
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
		return utils.HandleTxError(tx, err, LogTag)
	}

	scenarioInsertStmt, err := common.GetStmtRawQuery(tx, InsertScenarioTableStatement)

	if err != nil {
		zkLogger.Error(LogTag, "Error while creating the scenario insert scenarioVersionStmt ", err)
		return utils.HandleTxError(tx, err, LogTag)
	}

	params := []any{clusterId, request.ScenarioTitle, request.ScenarioType, time.Now().Unix()}
	scenarioId := 1000

	insertErr := zkPostgresRepo.dbRepo.InsertWithReturnRow(scenarioInsertStmt, params, []any{&scenarioId})

	if insertErr != nil {
		zkLogger.Error(LogTag, "Error while executing the insert query ", err)
		return utils.HandleTxError(tx, err, LogTag)
	}

	zkLogger.Debug(LogTag, "New scenarioId is ", scenarioId)

	scenarioObj := request.CreateScenarioObj(scenarioId)
	scenarioData, err := json.Marshal(scenarioObj)

	if err != nil {
		zkLogger.Error(LogTag, "Error while serializing scenario data ", err)
		return utils.HandleTxError(tx, err, LogTag)
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

	scenarioVersionStmt, err := common.GetStmtRawQuery(tx, InsertScenarioVersionTableStatement)

	if err != nil {
		zkLogger.Error(LogTag, "Error while creating the scenario version insert scenarioVersionStmt ", err)
		return utils.HandleTxError(tx, err, LogTag)
	}

	_, err = zkPostgresRepo.dbRepo.Insert(scenarioVersionStmt, scenarioVersionParams)

	if err != nil {
		zkLogger.Error(LogTag, "Error while inserting into the scenario version table. ", err)
		return utils.HandleTxError(tx, err, LogTag)
	}

	done, err2 := common.CommitTransaction(tx, LogTag)
	if err2 != nil {
		zkLogger.Error(LogTag, "Error while committing a db transaction in createNewScenario ", err2.Error)
		err = errors.New(err2.Error.Message)
		return err
	}

	if !done {
		zkLogger.Error(LogTag, "Transaction commit failed. ")
		return errors.New("Transaction commit failed. ")
	}

	zkLogger.Debug(LogTag, "Reached the end of the handler method.")
	return nil
}

func (zkPostgresRepo zkPostgresRepo) GetScenarioById(clusterId, scenarioId string) (*[]scenarioResponseModel.ScenarioDbResponse, error) {
	rows, err, closeRow := zkPostgresRepo.dbRepo.GetAll(GetScenarioByIdSqlStatement, []any{scenarioId, clusterId, DefaultClusterId})
	return Processor(rows, err, closeRow)
}

func (zkPostgresRepo zkPostgresRepo) GetAllScenario(filters *ScenarioQueryFilter) (*[]scenarioResponseModel.ScenarioDbResponse, error) {
	var params []any
	var query string
	if filters.Deleted == nil {
		query = GetAllScenarioSqlStatement
		params = []any{filters.Version, filters.ClusterId, DefaultClusterId, filters.Limit, filters.Offset}
	} else {
		query = GetAllScenarioForDashboardSqlStatement
		params = []any{filters.Version, filters.Deleted, filters.ClusterId, DefaultClusterId, filters.Limit, filters.Offset}
	}
	rows, err, closeRow := zkPostgresRepo.dbRepo.GetAll(query, params)

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
		err := rows.Scan(&scenarioResponse.ScenarioData, &scenarioResponse.Deleted, &scenarioResponse.Disabled, &scenarioResponse.CreatedAt, &scenarioResponse.DisabledAt, &scenarioResponse.UpdatedAt)
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

func (zkPostgresRepo zkPostgresRepo) DisableScenario(clusterId, scenarioId string, disable bool, disabledAtTime *int64, currentTime int64) (int, error) {
	count, err := updateScenario(zkPostgresRepo.dbRepo, DisableScenarioStatement, []any{disable, disabledAtTime, currentTime, clusterId, scenarioId})
	if err != nil {
		zkLogger.Error(LogTag, "Error in disable scenario ", err)
	}

	return count, err
}

func (zkPostgresRepo zkPostgresRepo) DeleteScenario(clusterId string, currentTime int64, scenarioId string) (int, error) {
	count, err := updateScenario(zkPostgresRepo.dbRepo, DeleteScenarioStatement, []any{currentTime, currentTime, clusterId, scenarioId})
	if err != nil {
		zkLogger.Error(LogTag, "Error in delete scenario ", err)
	}

	return count, err
}

func (zkPostgresRepo zkPostgresRepo) GetTotalRowsCount(filters *ScenarioQueryFilter) (int, error) {
	var count int
	params := []any{filters.Version, filters.ClusterId, DefaultClusterId}
	err := zkPostgresRepo.dbRepo.Get(GetTotalRowsCountStatement, params, []any{&count})
	if err != nil {
		zkLogger.Error(LogTag, "Error in GetTotalRowsCount ", err)
		return 0, err
	}

	return count, nil
}

func updateScenario(repo sqlDB.DatabaseRepo, query string, params []any) (int, error) {
	tx, err := repo.CreateTransaction()
	if err != nil {
		zkLogger.Info(LogTag, "Error Creating transaction")
		return 0, err
	}

	stmt, err := common.GetStmtRawQuery(tx, query)
	if err != nil {
		zkLogger.Info(LogTag, "Error Creating statement for disable/delete scenario", err)
		return 0, err
	}

	results, err := repo.Update(stmt, params)
	if err != nil {
		zkLogger.Info(LogTag, "Error in disable scenario ", err)
		return 0, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		zkLogger.Info(LogTag, "Error in disable/delete scenario ", err)
		return 0, err
	}

	zkLogger.Info(LogTag, "disable/disable count:", rowsAffected)
	if rowsAffected > 1 {
		zkLogger.Error(LogTag, "More than one scenario disabled/deleted, ROLLING BACK")
		done, err := common.RollbackTransaction(tx, LogTag)
		if err != nil {
			zkLogger.Error(LogTag, "Error while rolling back a db transaction in disable/delete Scenario ", err.Error)
			return 0, errors.New("Error while rolling back a db transaction in disable/delete Scenario ")
		}

		if done {
			return 0, nil
		}
	}

	_, zkErr := common.CommitTransaction(tx, LogTag)
	if zkErr != nil {
		zkLogger.Error(LogTag, "Error while committing a db transaction in disable/delete Scenario ", err.Error)
		return 0, errors.New("Error while committing a db transaction in disable/delete Scenario ")
	}

	return int(rowsAffected), nil
}
