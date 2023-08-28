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
)

const (
	DefaultClusterId                    = "Zk_default_cluster_id_for_all_scenarios"
	GetAllScenarioSqlStatement          = `SELECT scenario_data, deleted, disabled, created_at, disabled_at FROM scenario s INNER JOIN scenario_version sv USING(scenario_id) WHERE (scenario_version>$1 OR deleted_at>$2 OR disabled_at>$3) AND (cluster_id=$4 OR cluster_id=$5)`
	InsertScenarioTableStatement        = "INSERT INTO scenario (cluster_id, scenario_title, scenario_type) VALUES ($1, $2, $3) RETURNING scenario_id"
	InsertScenarioVersionTableStatement = "INSERT INTO scenario_version (scenario_id, scenario_data, schema_version, scenario_version, created_by, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	DisableScenarioStatement            = "UPDATE scenario set disabled=TRUE, disabled_at=$1 where cluster_id= $2 AND scenario_id=$3"
	DeleteScenarioStatement             = "UPDATE scenario set deleted=TRUE, deleted_at=$1 where cluster_id= $2 AND scenario_id=$3"
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
	DisableScenario(clusterId, scenarioId string) (int, error)
	DeleteScenario(clusterId, scenarioId string) (int, error)
}

type zkPostgresRepo struct {
	dbRepo sqlDB.DatabaseRepo
}

func NewZkPostgresRepo(db sqlDB.DatabaseRepo) ScenarioRepo {
	return &zkPostgresRepo{db}
}

func handleTxError(tx *sql.Tx, err2 error) error {
	done, err := common.RollbackTransaction(tx, LogTag)
	if err != nil {
		zkLogger.Error(LogTag, "Error while rolling back the transaction ", err.Error)
	}
	if !done {
		zkLogger.Error(LogTag, "Rolling back the transaction failed.")
	}
	return err2
}

func (zkPostgresRepo zkPostgresRepo) CreateNewScenario(clusterId string, request scenarioResponseModel.CreateScenarioRequest) error {
	tx, err := zkPostgresRepo.dbRepo.CreateTransaction()

	if err != nil {
		zkLogger.Error(LogTag, "Error while creating a db transaction in createNewScenario ", err)
		return handleTxError(tx, err)
	}

	scenarioInsertStmt, err := common.GetStmtRawQuery(tx, InsertScenarioTableStatement)

	if err != nil {
		zkLogger.Error(LogTag, "Error while creating the scenario insert scenarioVersionStmt ", err)
		return handleTxError(tx, err)
	}

	params := []any{clusterId, request.ScenarioTitle, request.ScenarioType}
	scenarioId := 1000

	insertErr := zkPostgresRepo.dbRepo.InsertWithReturnRow(scenarioInsertStmt, params, []any{&scenarioId})

	if insertErr != nil {
		zkLogger.Error(LogTag, "Error while executing the insert query ", err)
		return handleTxError(tx, err)
	}

	zkLogger.Debug(LogTag, "New scenarioId is ", scenarioId)

	scenarioObj := request.CreateScenarioObj(scenarioId)
	scenarioData, err := json.Marshal(scenarioObj)

	if err != nil {
		zkLogger.Error(LogTag, "Error while serializing scenario data ", err)
		return handleTxError(tx, err)
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
		return handleTxError(tx, err)
	}

	_, err = zkPostgresRepo.dbRepo.Insert(scenarioVersionStmt, scenarioVersionParams)

	if err != nil {
		zkLogger.Error(LogTag, "Error while inserting into the scenario version table. ", err)
		return handleTxError(tx, err)
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
		err := rows.Scan(&scenarioResponse.ScenarioData, &scenarioResponse.Deleted, &scenarioResponse.Disabled, &scenarioResponse.CreatedAt, &scenarioResponse.DisabledAt)
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

func (zkPostgresRepo zkPostgresRepo) DisableScenario(clusterId, scenarioId string) (int, error) {
	count, err := updateScenario(zkPostgresRepo.dbRepo, DisableScenarioStatement, clusterId, scenarioId, []any{time.Now().Unix(), clusterId, scenarioId})
	if err != nil {
		zkLogger.Error(LogTag, "Error in disable scenario ", err)
	}

	return count, err
}

func (zkPostgresRepo zkPostgresRepo) DeleteScenario(clusterId, scenarioId string) (int, error) {
	count, err := updateScenario(zkPostgresRepo.dbRepo, DeleteScenarioStatement, clusterId, scenarioId, []any{time.Now().Unix(), clusterId, scenarioId})
	if err != nil {
		zkLogger.Error(LogTag, "Error in delete scenario ", err)
	}

	return count, err
}

func updateScenario(repo sqlDB.DatabaseRepo, query, clusterId, scenarioId string, params []any) (int, error) {
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

	results, err := repo.Update(stmt, []any{time.Now().Unix(), clusterId, scenarioId})
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

	return int(rowsAffected), nil
}
