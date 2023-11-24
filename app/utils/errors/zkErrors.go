package errors

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
)

var (
	ZkErrorBadRequestTimeFormat                      = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Time not in supported format"}
	ZkErrorBadRequestServiceNameEmpty                = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Service name cannot be empty"}
	ZkErrorBadRequestServicePodEmpty                 = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Pod name cannot be empty"}
	ZkErrorBadRequestNamespaceEmpty                  = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Namespace cannot be empty"}
	ZkErrorBadRequestTimeEmpty                       = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Time cannot be empty"}
	ZkErrorBadRequestFileAttachedError               = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Error is the file attached"}
	ZkErrorBadRequestZkApiKeyEmpty                   = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Api Key cannot be empty"}
	ZkErrorBadRequestZkApiKeyMiddlewareEmpty         = zkerrors.ZkErrorType{Status: iris.StatusUnauthorized, Type: "BAD_REQUEST", Message: "Api Key cannot be empty"}
	ZkErrorBadRequestClusterIdEmpty                  = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "ClusterId cannot be empty"}
	ZkErrorBadRequestActionEmpty                     = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Action cannot be empty"}
	ZkErrorBadRequestActionInvalid                   = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Action can only be enable or disable"}
	ZkErrorBadRequestScenarioIdEmpty                 = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Scenario cannot be empty"}
	ZkErrorBadRequestVersionIsNotInteger             = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Filter Version is not integer"}
	ZkErrorBadRequestVersionEmpty                    = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Filter Version is empty"}
	ZkErrorBadRequestDeletedIsNotBoolean             = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "deleted is not bool"}
	ZkErrorBadRequestInvalidClusterAndUrlCombination = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "This integration does not exist for this cluster"}
	ZkErrorBadRequestUrl                             = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "url cannot be empty"}

	ZkErrorBadRequestDifferentVersions      = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Different versions in sheet"}
	ZkErrorBadRequestEmptyOrInvalidVersions = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Empty or Invalid versions in sheet"}
	ZkErrorBadRequestEmptyAttributeId       = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Empty attribute id in sheet"}
	ZkErrorBadRequestEmptyAttributePath     = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Empty attribute path in sheet"}
	ZkErrorBadRequestEmptyDataType          = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Empty data type in sheet"}
	ZkErrorBadRequestEmptyValue             = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Empty value in sheet"}
	ZkErrorBadRequestEmptyInput             = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Empty input in sheet"}
	ZkErrorBadRequestEmptyField             = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Empty field in sheet"}
	ZkErrorBadRequestEmptyProtocol          = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Empty protocol in sheet"}
	ZkErrorBadRequestEmptyKeySetName        = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Empty key set name in sheet"}
	ZkErrorBadRequestEmptyExecutor          = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Empty executor in sheet"}
	ZkErrorBadRequestSendToFrontend         = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Invalid value for send_to_frontend in sheet"}
	ZkErrorBadRequestJSON                   = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Invalid value for JSON in sheet"}
	ZkErrorBadRequestErrorInReadingFile     = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Error in reading file"}
	ZkErrorBadRequestFileNotFound           = zkerrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "File not found"}
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
