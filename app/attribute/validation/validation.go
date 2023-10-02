package validation

import (
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/logs"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"zk-api-server/app/attribute/model"
	"zk-api-server/app/utils"
	"zk-api-server/app/utils/errors"
)

var LogTag = "attribute_validation"

func ValidateGetAttributes(protocol string) *zkerrors.ZkError {
	if common.IsEmpty(protocol) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestVersionEmpty, nil)
		return &zkErr
	}
	return nil
}

func IsAttributesListValid(attributesList []model.AttributeInfoRequest) (bool, *zkerrors.ZkError) {
	if len(attributesList) == 0 {
		logger.Error(LogTag, "attributes list empty")
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, nil)
		return false, &zkError
	}

	if valid, zkErr := IsAllVersionSame(attributesList); !valid {
		return false, zkErr
	}

	if attributesList[0].Version == "common" {
		return ValidateCommonAttributesList(attributesList)
	} else {
		return ValidateVersionSpecificAttributesList(attributesList)
	}
}

func IsAllVersionSame(attributesList []model.AttributeInfoRequest) (bool, *zkerrors.ZkError) {
	version := attributesList[0].Version
	for _, v := range attributesList {
		if v.Version != version {
			zkLogger.Error(LogTag, "version mismatch")
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestDifferentVersions, nil)
			return false, &zkError
		}
	}
	return true, nil
}

func ValidateCommonAttributesList(attributesList []model.AttributeInfoRequest) (bool, *zkerrors.ZkError) {
	for i, v := range attributesList {
		if common.IsEmpty(v.Version) || v.Version != "common" {
			zkLogger.Error(LogTag, "version is empty or invalid at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyOrInvalidVersions, nil)
			return false, &zkError
		}

		if common.IsEmpty(v.AttributeId) {
			zkLogger.Error(LogTag, "attribute_id is empty at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyAttributeId, nil)
			return false, &zkError
		}

		if common.IsEmpty(v.AttributePath) {
			zkLogger.Error(LogTag, "attribute_path is empty at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyAttributePath, nil)
			return false, &zkError
		}

		if v.Field == nil || common.IsEmpty(*v.Field) {
			zkLogger.Error(LogTag, "field is empty at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyField, nil)
			return false, &zkError
		}

		if v.DataType == nil || common.IsEmpty(*v.DataType) {
			zkLogger.Error(LogTag, "data_type is empty at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyDataType, nil)
			return false, &zkError
		}

		if v.Input == nil || common.IsEmpty(*v.Input) {
			zkLogger.Error(LogTag, "input is empty at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyInput, nil)
			return false, &zkError
		}

		if *v.Input == "select" && (v.Values == nil || common.IsEmpty(*v.Values)) {
			zkLogger.Error(LogTag, "values is empty at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyValue, nil)
			return false, &zkError
		}

		if common.IsEmpty(v.Protocol) || (v.Protocol != "HTTP" && v.Protocol != "GENERAL") {
			zkLogger.Error(LogTag, "protocol is empty or invalid at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyProtocol, nil)
			return false, &zkError
		}

		if v.KeySetName == nil || common.IsEmpty(*v.KeySetName) {
			zkLogger.Error(LogTag, "key_set_name is empty at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyKeySetName, nil)
			return false, &zkError
		}

		if common.IsEmpty(v.Executor) || (v.Executor != utils.OTEL && v.Executor != utils.EBPF) {
			zkLogger.Error(LogTag, "executor is empty or invalid at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyExecutor, nil)
			return false, &zkError
		}
	}

	return true, nil
}

func ValidateVersionSpecificAttributesList(attributesList []model.AttributeInfoRequest) (bool, *zkerrors.ZkError) {
	for i, v := range attributesList {
		if common.IsEmpty(v.Version) || v.Version == "common" {
			zkLogger.Error(LogTag, "version is empty or invalid for version specific attributes list at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyOrInvalidVersions, nil)
			return false, &zkError
		}

		if common.IsEmpty(v.AttributeId) {
			zkLogger.Error(LogTag, "attribute_id is empty at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyAttributeId, nil)
			return false, &zkError
		}

		if common.IsEmpty(v.AttributePath) {
			zkLogger.Error(LogTag, "attribute_path is empty at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyAttributePath, nil)
			return false, &zkError
		}

		if common.IsEmpty(v.Protocol) || (v.Protocol != "HTTP" && v.Protocol != "GENERAL") {
			zkLogger.Error(LogTag, "protocol is empty or invalid at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyProtocol, nil)
			return false, &zkError
		}

		if common.IsEmpty(v.Executor) || (v.Executor != utils.OTEL && v.Executor != utils.EBPF) {
			zkLogger.Error(LogTag, "executor is empty or invalid at line: ", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestEmptyExecutor, nil)
			return false, &zkError
		}
	}

	return true, nil
}
