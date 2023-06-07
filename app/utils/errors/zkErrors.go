package errors

import (
	"github.com/kataras/iris/v12"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
)

var (
	ZkErrorBadRequestTimeFormat                = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Time not in supported format"}
	ZkErrorBadRequestServiceNameEmpty          = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Service name cannot be empty"}
	ZkErrorBadRequestServicePodEmpty           = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Pod name cannot be empty"}
	ZkErrorBadRequestNamespaceEmpty            = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Namespace cannot be empty"}
	ZkErrorBadRequestTimeEmpty                 = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Time cannot be empty"}
	ZkErrorBadRequestZkApiKeyEmpty             = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Api Key cannot be empty"}
	ZkErrorBadRequestZkApiKeyMiddlewareEmpty   = zkerrors.ZkErrorType{Status: iris.StatusUnauthorized, Type: "BAD_REQUEST", Message: "Api Key cannot be empty"}
	ZkErrorBadRequestClusterIdEmpty            = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "ClusterId cannot be empty"}
	ZkErrorBadRequestVersionIsNotInteger       = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Filter Version is not integer"}
	ZkErrorBadRequestVersionEmpty              = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Filter Version is not an integer"}
	ZkErrorBadRequestDeletedIsNotBoolean       = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "deleted is not bool"}
	ZkErrorBadRequestDefaultFilterIsNotBoolean = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "default_filter is not bool"}
)
