package zkerrors

import (
	"github.com/kataras/iris/v12"
)

type ZkErrorType struct {
	Message  string `json:"message"`
	Type     string `json:"type"`
	Status   int    `json:"status"`
	Metadata any    `json:"metadata"`
}

type ZkError struct {
	Error    ZkErrorType `json:"error"`
	Metadata any         `json:"metadata"`
}

// type _zkErrorBuilder struct {
// }

type ZkErrorBuilder struct {
}

var (
	ZK_ERROR_INTERNAL_SERVER                = ZkErrorType{Status: iris.StatusInternalServerError, Type: "INTERNAL_SERVER_ERROR", Message: "Encountered an issue, contact support"}
	ZK_ERROR_BAD_REQUEST_TIME_FORMAT        = ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Time not in supported format"}
	ZK_ERROR_TIMEOUT                        = ZkErrorType{Status: iris.StatusInternalServerError, Type: "OPERATION_TIMEOUT", Message: "Encountered an issue, contact support"}
	ZK_ERROR_NOT_FOUND                      = ZkErrorType{Status: iris.StatusNotFound, Type: "ENITITY_NOT_FOUND", Message: "Encountered an issue, contact support"}
	ZK_ERROR_SESSION_EXPIRED                = ZkErrorType{Status: iris.StatusPageExpired, Type: "SESSION_EXPIRED", Message: "The session has expired. Please login again to continue"}
	ZK_ERROR_UNAUTHORIZED                   = ZkErrorType{Status: iris.StatusUnauthorized, Type: "UNAUTHORIZED", Message: "You are unauthorized to perform this operation. Contact system admin"}
	ZK_ERROR_BAD_REQUEST                    = ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Encountered an issue while processing your request. Please check and try again!"}
	ZK_ERROR_INTERNAL_SERVER_SERVER         = ZkErrorType{Status: iris.StatusInternalServerError, Type: "INTERNAL_SERVER_ERROR", Message: "Encountered an issue while sending email, contact support"}
	ZK_ERROR_BAD_REQUEST_SERVICE_NAME_EMPTY = ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Service name cannot be empty"}
	ZK_ERROR_BAD_REQUEST_ZK_API_KEY_EMPTY   = ZkErrorType{Status: iris.StatusUnauthorized, Type: "BAD_REQUEST", Message: "API Key cannot be empty"}
	ZK_ERROR_BAD_REQUEST_SERVICE_POD_EMPTY  = ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Pod name cannot be empty"}
	ZK_ERROR_BAD_REQUEST_NAMESPACE_EMPTY    = ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Namespace cannot be empty"}
	ZK_ERROR_BAD_REQUEST_TIME_EMPTY         = ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Time cannot be empty"}
	ZK_ERROR_BAD_REQUEST_API_KEY_EMPTY      = ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Api Key cannot be empty"}
	ZK_ERROR_BAD_REQUEST_CLUSTER_ID_EMPTY   = ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "ClusterId cannot be empty"}
)

func (zkError ZkError) SetMetadata(metadata any) {
	zkError.Metadata = metadata
}

func (zkErrorBuilder ZkErrorBuilder) Build(zkErrorType ZkErrorType, metadata any) ZkError {
	return ZkError{
		Error:    zkErrorType,
		Metadata: metadata,
	}
}
