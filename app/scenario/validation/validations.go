package validation

import (
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
	"zk-api-server/app/scenario/model"
	"zk-api-server/app/utils"
	"zk-api-server/app/utils/errors"
)

func ValidateDisableScenarioApi(clusterId, scenarioId string, request model.ScenarioState) *zkerrors.ZkError {
	if zkCommon.IsEmpty(request.Action) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestActionEmpty, nil)
		return &zkErr
	}

	if request.Action != utils.Disable && request.Action != utils.Enable {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestActionInvalid, nil)
		return &zkErr
	}

	return validateScenarioIdAndClusterId(clusterId, scenarioId)
}

func ValidateGetScenarioByIdApi(clusterId, scenarioId string) *zkerrors.ZkError {
	return validateScenarioIdAndClusterId(clusterId, scenarioId)
}

func ValidateDeleteScenarioApi(clusterId, scenarioId string) *zkerrors.ZkError {
	return validateScenarioIdAndClusterId(clusterId, scenarioId)
}

func validateScenarioIdAndClusterId(clusterId, scenarioId string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(clusterId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestClusterIdEmpty, nil)
		return &zkErr
	}
	if zkCommon.IsEmpty(scenarioId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestScenarioIdEmpty, nil)
		return &zkErr
	}
	return nil
}

func ValidateGetAllScenarioApi(clusterId, version, deleted, offset, limit string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(clusterId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestClusterIdEmpty, nil)
		return &zkErr
	}
	if zkCommon.IsEmpty(version) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestVersionEmpty, nil)
		return &zkErr
	}

	_, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestVersionIsNotInteger, nil)
		return &zkErr
	}

	if !zkCommon.IsEmpty(deleted) {
		_, err = strconv.ParseBool(deleted)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestDeletedIsNotBoolean, nil)
			return &zkErr
		}
	}

	if !zkCommon.IsEmpty(limit) {
		_, err = strconv.Atoi(limit)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestLimitIsNotInteger, nil)
			return &zkErr
		}
	}

	if !zkCommon.IsEmpty(offset) {
		_, err = strconv.Atoi(offset)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestOffsetIsNotInteger, nil)
			return &zkErr
		}
	}

	return nil
}
