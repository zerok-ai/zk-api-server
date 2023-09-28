package repository

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/interfaces"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	"zk-api-server/app/attribute/model"
	"zk-api-server/utils"
)

type AttributeRepo interface {
	UpsertAttributes(dto model.AttributeDtoList) (bool, error)
	GetAttributes(protocol string) (model.AttributeDtoList, error)
	GetAttributesForBackend(version string) (model.AttributeDtoList, error)
}

const (
	UpsertAttributesQuery    = "INSERT INTO attributes (version, key_set, protocol, executor, attribute_list) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (version, key_set) DO UPDATE SET attribute_list = $5"
	SelectAttributesQuery    = "SELECT protocol, key_set, version, updated_at, attribute_list FROM attributes WHERE protocol=$1 AND version='common'"
	SelectAllAttributesQuery = "SELECT protocol, key_set, version, updated_at, attribute_list FROM attributes WHERE updated_at>$1  AND version!='common'"
)

var LogTag = "attributes_repo"

type zkPostgresRepo struct {
	dbRepo sqlDB.DatabaseRepo
}

func NewZkPostgresRepo(db sqlDB.DatabaseRepo) AttributeRepo {
	return &zkPostgresRepo{db}
}

func GetStmtRawQuery(tx *sql.Tx, stmt string) (*sql.Stmt, error) {
	preparedStmt, err := tx.Prepare(stmt)
	if err != nil {
		zkLogger.Error(LogTag, "Error preparing insert statement:", err)
		return nil, err
	}
	return preparedStmt, nil
}

func (z zkPostgresRepo) UpsertAttributes(dto model.AttributeDtoList) (bool, error) {
	tx, err := z.dbRepo.CreateTransaction()
	if err != nil {
		zkLogger.Info(LogTag, "Error Creating transaction")
		return false, err
	}

	attributeDetailsData := make([]interfaces.DbArgs, 0)

	for _, element := range dto {
		attributeDetailsData = append(attributeDetailsData, element)
	}

	stmt, err := GetStmtRawQuery(tx, UpsertAttributesQuery)
	if err != nil {
		zkLogger.Info(LogTag, "Error Creating statement for upsert", err)
		return false, err
	}

	dbRepo := z.dbRepo

	results, err := dbRepo.BulkUpsert(stmt, attributeDetailsData)
	if err != nil {
		zkLogger.Info(LogTag, "Error in bulk upsert ", err)
		return false, err
	}

	var upsertCount int64
	for _, v := range results {
		c, _ := v.RowsAffected()
		upsertCount += c
	}
	zkLogger.Info(LogTag, "bulk upsert count:", upsertCount)

	b, txErr := common.CommitTransaction(tx, LogTag)
	if txErr != nil || !b {
		zkLogger.Error(LogTag, "Error while committing the transaction ", txErr.Error)
		return false, utils.HandleTxError(tx, err, LogTag)
	}

	return true, nil
}

func (z zkPostgresRepo) GetAttributes(protocol string) (model.AttributeDtoList, error) {
	rows, err, closeRow := z.dbRepo.GetAll(SelectAttributesQuery, []any{protocol})
	return Processor(rows, err, closeRow)
}

func (z zkPostgresRepo) GetAttributesForBackend(updateAt string) (model.AttributeDtoList, error) {
	rows, err, closeRow := z.dbRepo.GetAll(SelectAllAttributesQuery, []any{updateAt})
	return Processor(rows, err, closeRow)
}

func Processor(rows *sql.Rows, sqlErr error, f func()) (model.AttributeDtoList, error) {
	defer f()

	if sqlErr != nil {
		return nil, sqlErr
	}

	if rows == nil {
		zkLogger.Debug(LogTag, "rows nil", sqlErr)
		return nil, sqlErr
	}

	var attributeDto model.AttributeDto
	var attributeDtoArr []model.AttributeDto
	for rows.Next() {
		err := rows.Scan(&attributeDto.Protocol, &attributeDto.KeySet, &attributeDto.Version, &attributeDto.UpdatedAt, &attributeDto.Attributes)
		if err != nil {
			zkLogger.Error(LogTag, err)
		}

		attributeDtoArr = append(attributeDtoArr, attributeDto)
	}

	err := rows.Err()
	if err != nil {
		zkLogger.Error(LogTag, err)
	}

	return attributeDtoArr, nil
}
