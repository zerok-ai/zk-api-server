package repository

import (
	"database/sql"
	"fmt"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	dto "zk-api-server/app/obfuscation/model/dto"
)

const (
	GetObfuscationById             = "SELECT id, org_id, rule_name, rule_type, rule_def, created_at, updated_at, deleted, disabled FROM zk_obfuscation WHERE id=$1 AND org_id=$2"
	GetAllObfuscationsOperator     = "SELECT id, org_id, rule_name, rule_type, rule_def, created_at, updated_at, deleted, disabled FROM zk_obfuscation WHERE org_id=$1 AND updated_at>$2"
	InsertObfuscation              = "INSERT INTO zk_obfuscation (org_id, rule_name, rule_type, rule_def, created_at, updated_at, deleted, disabled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	UpdateObfuscation              = "UPDATE zk_obfuscation SET rule_name=$1, rule_type=$2, rule_def=$3, updated_at=$4, deleted=$5, disabled=$6 WHERE id=$7 AND org_id=$8"
	DeleteObfuscation              = "UPDATE zk_obfuscation SET deleted = true WHERE id=$1 AND org_id=$2"
	GetAllObfuscationsForDashboard = "SELECT id, org_id, rule_name, rule_type, rule_def, created_at, updated_at, deleted, disabled FROM zk_obfuscation WHERE org_id=$1 AND deleted = false ORDER BY updated_at LIMIT $2 OFFSET $3"
	GetTotalRowsCountStatement     = "SELECT COUNT(*) as count FROM zk_obfuscation WHERE deleted=false AND org_id=$1"
)

var LogTag = "obfuscations_repo"

type ObfuscationRepo interface {
	GetAllObfuscationsForDashboard(orgId string, offset, limit string) ([]dto.Obfuscation, int, error)
	GetAllObfuscationsOperator(orgId string, updatedTime int64) ([]dto.Obfuscation, error)
	GetObfuscationById(id string, orgId string) (*dto.Obfuscation, error)
	InsertObfuscation(obfuscation dto.Obfuscation) (bool, error)
	UpdateObfuscation(obfuscation dto.Obfuscation) (bool, error)
	DeleteObfuscation(orgId string, id string) (bool, error)
}

// ObfuscationRepo implementation
type zkPostgresObfuscationRepo struct {
	dbRepo sqlDB.DatabaseRepo
}

func (z zkPostgresObfuscationRepo) GetAllObfuscationsOperator(orgId string, updatedTime int64) ([]dto.Obfuscation, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetAllObfuscationsOperator, []any{orgId, updatedTime})
	return ObfuscationProcessor(rows, err, closeRow)
}

func NewZkPostgresObfuscationRepo(db sqlDB.DatabaseRepo) ObfuscationRepo {
	return &zkPostgresObfuscationRepo{db}
}

func (z zkPostgresObfuscationRepo) GetAllObfuscationsForDashboard(orgId string, offset, limit string) ([]dto.Obfuscation, int, error) {
	// Creating a Repeatable read transaction here, since we are running two queries here:
	// One to get the obfuscations with limit and offset and another to get the total rows count.
	// In case of an update in the obfuscation table, the total rows count might change after the first query is executed.
	// This will lead to a wrong total rows count give to the dashboard.
	// So doing a Repeatable read transaction here to avoid this issue.
	tx, err := z.dbRepo.CreateTransactionWithIsolation(sql.LevelRepeatableRead)
	if err != nil {
		z.rollbackTransaction(tx)
		zkLogger.Error(LogTag, "Error while creating a transaction ", err)
		return nil, 0, err
	}
	// Function to handle panic and rollback transaction
	handlePanic := func(tx *sql.Tx) {
		if r := recover(); r != nil {
			zkLogger.Error(LogTag, "Recovered from panic: ", err)
			z.rollbackTransaction(tx)
		}
	}

	defer handlePanic(tx)

	// Get all obfuscations for dashboard
	rows, err, closeRow := z.dbRepo.GetAllWithTx(tx, GetAllObfuscationsForDashboard, []any{orgId, limit, offset})
	if err != nil {
		zkLogger.Error(LogTag, "Error while getting all obfuscations for dashboard ", err)
		z.rollbackTransaction(tx)
		return nil, 0, err
	}

	// Process obfuscations
	obfuscations, err := ObfuscationProcessor(rows, err, closeRow)
	if err != nil {
		zkLogger.Error(LogTag, "Error while processing obfuscations ", err)
		z.rollbackTransaction(tx)
		return nil, 0, err
	}

	// Get total rows count
	var count int
	params := []any{orgId}
	err = z.dbRepo.GetWithTx(tx, GetTotalRowsCountStatement, params, []any{&count})
	if err != nil {
		zkLogger.Error(LogTag, "Error in GetTotalRowsCount ", err)
		z.rollbackTransaction(tx)
		return nil, 0, err
	}

	z.dbRepo.CommitTransaction(tx)
	return obfuscations, count, nil
}

func (z zkPostgresObfuscationRepo) rollbackTransaction(tx *sql.Tx) error {
	err := z.dbRepo.RollbackTransaction(tx)
	if err != nil {
		zkLogger.Error(LogTag, "Error while rolling back a transaction ", err)
		return err
	}
	return nil
}

func (z zkPostgresObfuscationRepo) GetObfuscationById(id string, orgId string) (*dto.Obfuscation, error) {
	var row dto.Obfuscation
	err := z.dbRepo.Get(GetObfuscationById, []any{id, orgId}, []any{&row.ID, &row.OrgID, &row.RuleName, &row.RuleType, &row.RuleDef, &row.CreatedAt, &row.UpdatedAt, &row.Deleted, &row.Disabled})
	if err != nil {
		zkLogger.Error(LogTag, "Error while getting the obfuscation by id: ", id, err)
		return nil, err
	}
	return &row, nil
}

func (z zkPostgresObfuscationRepo) InsertObfuscation(obfuscation dto.Obfuscation) (bool, error) {
	if obfuscation.Deleted || obfuscation.Disabled {
		zkLogger.Error(LogTag, "Obfuscation cannot be deleted or disabled.")
		return false, fmt.Errorf("cannot insert deleted or disabled obfuscation")
	}
	obfuscationInsertStmt := z.dbRepo.CreateStatement(InsertObfuscation)

	result, err := z.dbRepo.Insert(obfuscationInsertStmt, obfuscation)
	if err != nil {
		zkLogger.Error(LogTag, err)
		return false, err
	}

	zkLogger.Info(LogTag, "Obfuscation insert successfully ", result.RowsAffected)

	return true, nil
}

func (z zkPostgresObfuscationRepo) UpdateObfuscation(obfuscation dto.Obfuscation) (bool, error) {
	obfuscationUpdateStmt := z.dbRepo.CreateStatement(UpdateObfuscation)

	values := []any{obfuscation.RuleName, obfuscation.RuleType, obfuscation.RuleDef, obfuscation.UpdatedAt, obfuscation.Deleted, obfuscation.Disabled, obfuscation.ID, obfuscation.OrgID}

	result, err := z.dbRepo.Update(obfuscationUpdateStmt, values)
	if err != nil {
		zkLogger.Error(LogTag, err)
		return false, err
	}

	zkLogger.Info(LogTag, "Obfuscation updated successfully ", result.RowsAffected)

	return true, nil
}

func (z zkPostgresObfuscationRepo) DeleteObfuscation(orgId string, id string) (bool, error) {
	obfuscationUpdateStmt := z.dbRepo.CreateStatement(DeleteObfuscation)

	values := []any{id, orgId}

	result, err := z.dbRepo.Update(obfuscationUpdateStmt, values)
	if err != nil {
		zkLogger.Error(LogTag, err)
		return false, err
	}

	zkLogger.Info(LogTag, "Obfuscation Deleted successfully ", result.RowsAffected)

	return true, nil
}

func ObfuscationProcessor(rows *sql.Rows, sqlErr error, f func()) ([]dto.Obfuscation, error) {
	defer f()

	if sqlErr != nil {
		return nil, sqlErr
	}

	if rows == nil {
		zkLogger.Debug(LogTag, "rows nil", sqlErr)
		return nil, sqlErr
	}

	var obfuscation dto.Obfuscation
	var obfuscationsArr []dto.Obfuscation
	for rows.Next() {
		err := rows.Scan(&obfuscation.ID, &obfuscation.OrgID, &obfuscation.RuleName, &obfuscation.RuleType, &obfuscation.RuleDef, &obfuscation.CreatedAt, &obfuscation.UpdatedAt, &obfuscation.Deleted, &obfuscation.Disabled)
		if err != nil {
			zkLogger.Error(LogTag, err)
		}

		obfuscationsArr = append(obfuscationsArr, obfuscation)
	}

	err := rows.Err()
	if err != nil {
		zkLogger.Error(LogTag, err)
	}

	return obfuscationsArr, nil
}
