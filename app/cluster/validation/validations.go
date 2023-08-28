package validation

import (
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
	"zk-api-server/app/utils"
	"zk-api-server/app/utils/errors"
)

func ValidatePxlTime(s string) bool {
	if !utils.IsValidPxlTime(s) {
		return false
	}
	return true
}

func ValidateGraphDetailsApi(serviceName, ns, st, apiKey string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(serviceName) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestServiceNameEmpty, nil)
		return &zkErr
	}
	if zkCommon.IsEmpty(ns) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestNamespaceEmpty, nil)
		return &zkErr
	}
	if zkCommon.IsEmpty(st) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeEmpty, nil)
		return &zkErr
	}
	if zkCommon.IsEmpty(apiKey) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestZkApiKeyEmpty, nil)
		return &zkErr
	}
	return nil
}

func ValidatePodDetailsApi(podName, ns, st, apiKey string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(podName) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestServicePodEmpty, nil)
		return &zkErr
	}
	if zkCommon.IsEmpty(ns) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestNamespaceEmpty, nil)
		return &zkErr
	}
	if zkCommon.IsEmpty(st) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeEmpty, nil)
		return &zkErr
	}
	if zkCommon.IsEmpty(apiKey) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestZkApiKeyEmpty, nil)
		return &zkErr
	}
	return nil
}

func ValidateGetResourceDetailsApi(st string, apiKey string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(st) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeEmpty, nil)
		return &zkErr
	}
	if zkCommon.IsEmpty(apiKey) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestZkApiKeyEmpty, nil)
		return &zkErr
	}
	return nil
}

func ValidateDisableScenarioApi(clusterId, scenarioId string) *zkerrors.ZkError {
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

func ValidateGetPxlData(s string, apiKey string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(s) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestClusterIdEmpty, nil)
		return &zkErr
	}
	if zkCommon.IsEmpty(apiKey) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestZkApiKeyEmpty, nil)
		return &zkErr
	}
	return nil
}
