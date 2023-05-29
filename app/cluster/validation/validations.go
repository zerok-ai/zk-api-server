package validation

import (
	"main/app/utils"
	"main/app/utils/zkerrors"
	"strconv"
)

func ValidatePxlTime(s string) bool {
	if !utils.IsValidPxlTime(s) {
		return false
	}
	return true
}

func ValidateGraphDetailsApi(serviceName, ns, st, apiKey string) *zkerrors.ZkError {
	if utils.IsEmpty(serviceName) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_SERVICE_NAME_EMPTY, nil)
		return &zkErr
	}
	if utils.IsEmpty(ns) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_NAMESPACE_EMPTY, nil)
		return &zkErr
	}
	if utils.IsEmpty(st) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_EMPTY, nil)
		return &zkErr
	}
	if utils.IsEmpty(apiKey) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_ZK_API_KEY_EMPTY, nil)
		return &zkErr
	}
	return nil
}

func ValidatePodDetailsApi(podName, ns, st, apiKey string) *zkerrors.ZkError {
	if utils.IsEmpty(podName) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_SERVICE_POD_EMPTY, nil)
		return &zkErr
	}
	if utils.IsEmpty(ns) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_NAMESPACE_EMPTY, nil)
		return &zkErr
	}
	if utils.IsEmpty(st) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_EMPTY, nil)
		return &zkErr
	}
	if utils.IsEmpty(apiKey) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_ZK_API_KEY_EMPTY, nil)
		return &zkErr
	}
	return nil
}

func ValidateGetResourceDetailsApi(st string, apiKey string) *zkerrors.ZkError {
	if utils.IsEmpty(st) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_EMPTY, nil)
		return &zkErr
	}
	if utils.IsEmpty(apiKey) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_ZK_API_KEY_EMPTY, nil)
		return &zkErr
	}
	return nil
}

func ValidateGetAllScenarioApi(clusterId, version, deleted, offset, limit string) *zkerrors.ZkError {
	if utils.IsEmpty(clusterId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_CLUSTER_ID_EMPTY, nil)
		return &zkErr
	}
	if utils.IsEmpty(version) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_VERSION_EMPTY, nil)
		return &zkErr
	}

	_, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_VERSION_IS_NOT_INTEGER, nil)
		return &zkErr
	}

	if !utils.IsEmpty(deleted) {
		_, err = strconv.ParseBool(deleted)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_DELETED_IS_NOT_BOOLEAN, nil)
			return &zkErr
		}
	}

	if !utils.IsEmpty(limit) {
		_, err = strconv.Atoi(limit)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_LIMIT_IS_NOT_INTEGER, nil)
			return &zkErr
		}
	}

	if !utils.IsEmpty(offset) {
		_, err = strconv.Atoi(offset)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_OFFSET_IS_NOT_INTEGER, nil)
			return &zkErr
		}
	}

	return nil
}

func ValidateGetPxlData(s string, apiKey string) *zkerrors.ZkError {
	if utils.IsEmpty(s) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_CLUSTER_ID_EMPTY, nil)
		return &zkErr
	}
	if utils.IsEmpty(apiKey) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_ZK_API_KEY_EMPTY, nil)
		return &zkErr
	}
	return nil
}
