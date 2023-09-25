package validation

import (
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"zk-api-server/app/utils/errors"
)

func ValidateGetAttributes(protocol string) *zkerrors.ZkError {
	if common.IsEmpty(protocol) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestVersionEmpty, nil)
		return &zkErr
	}
	return nil
}
