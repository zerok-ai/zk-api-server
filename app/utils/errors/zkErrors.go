package errors

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
)

var (
	ZkErrorBadRequestTimeFormat              = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Time not in supported format"}
	ZkErrorBadRequestServiceNameEmpty        = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Service name cannot be empty"}
	ZkErrorBadRequestServicePodEmpty         = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Pod name cannot be empty"}
	ZkErrorBadRequestNamespaceEmpty          = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Namespace cannot be empty"}
	ZkErrorBadRequestTimeEmpty               = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Time cannot be empty"}
	ZkErrorBadRequestFileAttachedError       = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Error is the file attached"}
	ZkErrorBadRequestZkApiKeyEmpty           = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Api Key cannot be empty"}
	ZkErrorBadRequestZkApiKeyMiddlewareEmpty = zkerrors.ZkErrorType{Status: iris.StatusUnauthorized, Type: "BAD_REQUEST", Message: "Api Key cannot be empty"}
	ZkErrorBadRequestClusterIdEmpty          = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "ClusterId cannot be empty"}
	ZkErrorBadRequestActionEmpty             = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Action cannot be empty"}
	ZkErrorBadRequestActionInvalid           = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Action can only be enable or disable"}
	ZkErrorBadRequestScenarioIdEmpty         = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Scenario cannot be empty"}
	ZkErrorBadRequestVersionIsNotInteger     = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Filter Version is not integer"}
	ZkErrorBadRequestVersionEmpty            = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Filter Version is empty"}
	ZkErrorBadRequestDeletedIsNotBoolean     = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "deleted is not bool"}
	ZkErrorBadRequestKeySetEmpty             = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "key_set is empty"}
)

var (
	ErrAuthenticationFailed                   = errors.New("rpc error: code = Internal desc = Auth middleware failed: failed to fetch token - unauthenticated")
	ErrClusterParsingFailed                   = errors.New("failed to parse cluster info")
	ErrClusterIdEmpty                         = errors.New("clusterId cannot be empty")
	ErrPxlStartTimeEmpty                      = errors.New("start Time st cannot be empty")
	ErrZkApiKeyEmpty                          = errors.New("Zk-Api-Key header cannot be empty")
	ErrNamespaceEmpty                         = errors.New("namespace ns cannot be empty")
	ErrServiceNameEmpty                       = errors.New("service name cannot be empty")
	ErrPodNameEmpty                           = errors.New("pod name cannot be empty")
	ErrInternalServerError                    = errors.New("something went wrong, please try again later")
	ErrInvalidRuleDoesNotContainZkRequestType = errors.New("this filter rule does not contain zk request type rule")
	ErrUnableToAccessFile                     = errors.New("cannot access given file")
)
