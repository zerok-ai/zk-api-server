package utils

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/logs"
	"mime/multipart"
)

func ParseCSV(file multipart.File) ([][]string, error) {
	r := csv.NewReader(bufio.NewReader(file))

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func HandleTxError(tx *sql.Tx, txnErr error, logTag string) error {
	done, err := common.RollbackTransaction(tx, logTag)
	if err != nil {
		logger.Error(logTag, "Error while rolling back the transaction ", err.Error)
	}
	if !done {
		logger.Error(logTag, "Rolling back the transaction failed.")
	}
	return txnErr
}
