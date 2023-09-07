package repository

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	"zk-api-server/app/integrations/model/dto"
)

const (
	GetIntegrationById           = "SELECT id, cluster_id, type, url, authentication, level, created_at, updated_at, deleted, disabled FROM integrations WHERE id=$1 AND cluster_id=$2"
	GetAllActiveIntegrations     = "SELECT id, type, url, authentication, level, created_at, updated_at, deleted, disabled FROM integrations WHERE cluster_id=$1 AND deleted = false AND disabled = false"
	GetAllNonDeletedIntegrations = "SELECT id, type, url, authentication, level, created_at, updated_at, deleted, disabled FROM integrations WHERE cluster_id=$1 AND deleted = false"
	UpsertIntegration            = "INSERT INTO integrations (id, cluster_id, type, url, authentication, level, created_at, updated_at, deleted, disabled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (id) DO UPDATE SET type = $3, url = $4, authentication = $5, level = $6, deleted = $9, disabled = $10"
)

type IntegrationRepo interface {
	GetAllIntegrations(clusterId string, onlyActive bool) ([]dto.Integration, error)
	GetIntegrationsById(id int, clusterId string) (*dto.Integration, error)
	UpsertIntegration(integration dto.Integration) (bool, error)
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

func (z zkPostgresRepo) GetIntegrationsById(id int, clusterId string) (*dto.Integration, error) {
	var row dto.Integration
	err := z.dbRepo.Get(GetIntegrationById, []any{id, clusterId}, []any{&row.ID, &row.ClusterId, &row.Type, &row.URL, &row.Authentication, &row.Level, &row.CreatedAt, &row.UpdatedAt, &row.Deleted, &row.Disabled})
	if err != nil {
		zkLogger.Error(LogTag, "Error while getting the integration by id ", err)
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
		err := rows.Scan(&integration.ID, &integration.Type, &integration.URL, &integration.Authentication, &integration.Level, &integration.CreatedAt, &integration.UpdatedAt, &integration.Deleted, &integration.Disabled)
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

func (z zkPostgresRepo) UpsertIntegration(integration dto.Integration) (bool, error) {
	tx, err := z.dbRepo.CreateTransaction()
	integrationUpsertStmt, err := common.GetStmtRawQuery(tx, UpsertIntegration)

	if err != nil {
		zkLogger.Error(LogTag, "Error while creating the integration insert  ", err)
		return false, err
	}

	result, err := z.dbRepo.Upsert(integrationUpsertStmt, integration)
	if err != nil {
		zkLogger.Error(LogTag, err)
		return false, err
	}

	zkLogger.Info(LogTag, "Integration upsert successfully ", result.RowsAffected)

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
