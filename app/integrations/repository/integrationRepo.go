package repository

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	"zk-api-server/app/integrations/model/dto"
)

const (
	GetIntegrationById           = "SELECT id, cluster_id, alias, type, url, authentication, level, created_at, updated_at, deleted, disabled, metric_server FROM zk_integrations WHERE id=$1 AND cluster_id=$2"
	GetAllActiveIntegrations     = "SELECT id, alias, type, url, authentication, level, created_at, updated_at, deleted, disabled, metric_server FROM zk_integrations WHERE cluster_id=$1 AND deleted = false AND disabled = false"
	GetAnIntegrationDetails      = "SELECT id, cluster_id, alias, type, url, authentication, level, created_at, updated_at, deleted, disabled, metric_server FROM zk_integrations WHERE id=$1 AND deleted = false"
	GetAllNonDeletedIntegrations = "SELECT id, alias, type, url, authentication, level, created_at, updated_at, deleted, disabled, metric_server FROM zk_integrations WHERE cluster_id=$1 AND deleted = false"
	InsertIntegration            = "INSERT INTO zk_integrations (cluster_id, alias, type, url, authentication, level, created_at, updated_at, deleted, disabled, metric_server) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id"
	UpdateIntegration            = "UPDATE zk_integrations SET alias=$1, type = $2, url = $3, authentication = $4, level = $5, deleted = $6, disabled = $7, updated_at = $8, metric_server = $9 WHERE id = $10"
)

type IntegrationRepo interface {
	GetAllIntegrations(clusterId string, onlyActive bool) ([]dto.Integration, error)
	GetIntegrationsById(id string, clusterId string) (*dto.Integration, error)
	InsertIntegration(integration dto.Integration) (bool, int, error)
	UpdateIntegration(integration dto.Integration) (bool, error)
	GetAnIntegrationDetails(integrationId string) ([]dto.Integration, error)
}

var LogTag = "integrations_repo"

type zkPostgresRepo struct {
	dbRepo sqlDB.DatabaseRepo
}

func NewZkPostgresRepo(db sqlDB.DatabaseRepo) IntegrationRepo {
	return &zkPostgresRepo{db}
}

func (z zkPostgresRepo) GetAllIntegrations(clusterId string, onlyActive bool) ([]dto.Integration, error) {
	var query string
	if onlyActive {
		query = GetAllActiveIntegrations
	} else {
		query = GetAllNonDeletedIntegrations
	}

	rows, err, closeRow := z.dbRepo.GetAll(query, []any{clusterId})
	return Processor(rows, err, closeRow)
}

func (z zkPostgresRepo) GetAnIntegrationDetails(integrationId string) ([]dto.Integration, error) {
	query := GetAnIntegrationDetails
	rows, err, closeRow := z.dbRepo.GetAll(query, []any{integrationId})
	return Processor(rows, err, closeRow)
}

func (z zkPostgresRepo) GetIntegrationsById(id string, clusterId string) (*dto.Integration, error) {
	var row dto.Integration
	err := z.dbRepo.Get(GetIntegrationById, []any{id, clusterId}, []any{&row.ID, &row.ClusterId, &row.Alias, &row.Type, &row.URL, &row.Authentication, &row.Level, &row.CreatedAt, &row.UpdatedAt, &row.Deleted, &row.Disabled, &row.MetricServer})
	if err != nil {
		zkLogger.Error(LogTag, "Error while getting the integration by id: ", id, err)
		return nil, err
	}
	return &row, err
}

func Processor(rows *sql.Rows, sqlErr error, f func()) ([]dto.Integration, error) {
	defer f()

	if sqlErr != nil {
		return nil, sqlErr
	}

	if rows == nil {
		zkLogger.Debug(LogTag, "rows nil", sqlErr)
		return nil, sqlErr
	}

	var integration dto.Integration
	var integrationArr []dto.Integration
	for rows.Next() {
		err := rows.Scan(&integration.ID, &integration.ClusterId, &integration.Alias, &integration.Type, &integration.URL, &integration.Authentication, &integration.Level, &integration.CreatedAt, &integration.UpdatedAt, &integration.Deleted, &integration.Disabled, &integration.MetricServer)
		if err != nil {
			zkLogger.Error(LogTag, err)
		}

		integrationArr = append(integrationArr, integration)
	}

	err := rows.Err()
	if err != nil {
		zkLogger.Error(LogTag, err)
	}

	return integrationArr, nil
}

func (z zkPostgresRepo) InsertIntegration(integration dto.Integration) (bool, int, error) {
	tx, err := z.dbRepo.CreateTransaction()
	integrationUpsertStmt, err := common.GetStmtRawQuery(tx, InsertIntegration)

	integrationId := -1
	if err != nil {
		zkLogger.Error(LogTag, "Error while creating the integration insert  ", err)
		return false, integrationId, handleTxError(tx, err)
	}

	err = z.dbRepo.InsertWithReturnRow(integrationUpsertStmt, []any{integration}, []any{&integrationId})
	if err != nil {
		zkLogger.Error(LogTag, err)
		return false, integrationId, handleTxError(tx, err)
	}

	b, txErr := common.CommitTransaction(tx, LogTag)
	if txErr != nil || !b {
		zkLogger.Error(LogTag, "Error while committing the transaction ", txErr.Error)
		return false, integrationId, handleTxError(tx, err)
	}

	return true, integrationId, nil
}

func (z zkPostgresRepo) UpdateIntegration(integration dto.Integration) (bool, error) {
	tx, err := z.dbRepo.CreateTransaction()
	integrationUpsertStmt, err := common.GetStmtRawQuery(tx, UpdateIntegration)

	if err != nil {
		zkLogger.Error(LogTag, "Error while creating the integration update  ", err)
		return false, handleTxError(tx, err)
	}

	result, err := z.dbRepo.Update(integrationUpsertStmt, []any{integration.Alias, integration.Type, integration.URL, integration.Authentication, integration.Level, integration.Deleted, integration.Disabled, integration.UpdatedAt, integration.MetricServer, integration.ID})
	if err != nil {
		zkLogger.Error(LogTag, err)
		return false, handleTxError(tx, err)
	}

	b, txErr := common.CommitTransaction(tx, LogTag)
	if txErr != nil || !b {
		zkLogger.Error(LogTag, "Error while committing the transaction ", txErr.Error)
		return false, handleTxError(tx, err)
	}
	zkLogger.Info(LogTag, "Integration update successfully ", result.RowsAffected)

	return true, nil
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
