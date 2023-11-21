package utils

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/logs"
)

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
