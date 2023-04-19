package validation

import (
	"github.com/kataras/iris/v12"
	"main/app/utils"
	"main/app/utils/zkerrors"
)

func ValidatePxlTime(s string) bool {
	if !utils.IsValidPxlTime(s) {
		return false
	}
	return true
}

func ValidateGraphDetailsApi(ctx iris.Context, serviceName, ns, st, apiKey string) *zkerrors.ZkError {
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
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_API_KEY_EMPTY, nil)
		return &zkErr
	}
	return nil
}

func ValidatePodDetailsApi(ctx iris.Context, podName, ns, st, apiKey string) *zkerrors.ZkError {
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
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_API_KEY_EMPTY, nil)
		return &zkErr
	}
	return nil
}

func ValidateGetResourceDetailsApi(ctx iris.Context, st string, apiKey string) *zkerrors.ZkError {
	if utils.IsEmpty(st) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_EMPTY, nil)
		return &zkErr
	}
	if utils.IsEmpty(apiKey) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_API_KEY_EMPTY, nil)
		return &zkErr
	}
	return nil
}

func ValidateGetPxlData(ctx iris.Context, s string, apiKey string) *zkerrors.ZkError {
	if utils.IsEmpty(s) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_CLUSTER_ID_EMPTY, nil)
		return &zkErr
	}
	if utils.IsEmpty(apiKey) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_API_KEY_EMPTY, nil)
		return &zkErr
	}
	return nil
}
