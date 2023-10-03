package service

import (
	"encoding/csv"
	"fmt"
	"github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"mime/multipart"
	"strconv"
	"strings"
	"zk-api-server/app/attribute/model"
	"zk-api-server/app/attribute/repository"
	"zk-api-server/app/attribute/validation"
	"zk-api-server/app/utils/errors"
)

type AttributeService interface {
	GetAttributes(protocols []string) (*model.AttributeListResponse, *zkerrors.ZkError)
	GetAttributesForBackend(protocol string) (*model.ExecutorAttributesResponse, *zkerrors.ZkError)
	UpsertAttributes(multipart.File) (bool, *zkerrors.ZkError)
}

var LogTag = "attribute_service"

type attributeService struct {
	repo repository.AttributeRepo
}

func NewAttributeService(repo repository.AttributeRepo) AttributeService {
	return &attributeService{repo: repo}
}

func (a attributeService) GetAttributes(protocols []string) (*model.AttributeListResponse, *zkerrors.ZkError) {
	if len(protocols) == 0 {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, nil)
		zkLogger.Error(LogTag, "protocol is empty")
		return nil, &zkError
	}

	sanitizedProtocols := make([]string, 0)
	for _, protocol := range protocols {
		s := strings.Trim(protocol, " ")
		sanitizedProtocols = append(sanitizedProtocols, s)
	}

	data, err := a.repo.GetAttributes(sanitizedProtocols)
	if err != nil {
		zkLogger.Error(LogTag, "failed to get attributes list", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorDbError, err)
		return nil, &zkError

	}
	response := model.ConvertAttributeDtoToAttributeResponse(data)
	return &response, nil
}

func (a attributeService) GetAttributesForBackend(updatedAt string) (*model.ExecutorAttributesResponse, *zkerrors.ZkError) {
	updatedAtInt, _ := strconv.ParseInt(updatedAt, 10, 64)
	if common.IsEmpty(updatedAt) {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, nil)
		zkLogger.Error(LogTag, "protocol is empty")
		return nil, &zkError
	}

	data, err := a.repo.GetAttributesForBackend(updatedAt)
	if err != nil {
		zkLogger.Error(LogTag, "failed to get attributes list", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorDbError, err)
		return nil, &zkError

	}
	response := model.ConvertAttributeDtoToExecutorAttributesResponse(data)
	if response.Version > updatedAtInt {
		response.Update = true
	}

	return &response, nil
}

func readCSVAndReturnData(file multipart.File) ([]model.AttributeInfoRequest, *zkerrors.ZkError) {
	dtoList := make([]model.AttributeInfoRequest, 0)
	if file == nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestFileNotFound, nil)
		zkLogger.Error(LogTag, "file is nil")
		return dtoList, &zkError
	}

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading header:", err)
		zkLogger.Error(LogTag, "Error reading header:", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestErrorInReadingFile, nil)
		return dtoList, &zkError
	}

	colIndex := make(map[string]int)
	for i, colName := range header {
		colIndex[colName] = i
	}
	headersMap := colIndex

	rowCount := 0
	for {
		rowCount++
		row, err := reader.Read()
		if err != nil {
			break
		}

		sendToFrontEnd, err := strconv.ParseBool(row[headersMap["Send to Frontend"]])
		if err != nil {
			zkLogger.ErrorF(LogTag, "Error parsing CSV: %v at sendToFrontEnd, at line: %v", err, rowCount)
			zkError := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestSendToFrontend, nil)
			return dtoList, &zkError
		}

		if !sendToFrontEnd {
			continue
		}

		var supportedFormatsValue *[]string
		supportedFormatsStr := row[headersMap["supported_formats"]]
		if common.IsEmpty(supportedFormatsStr) {
			supportedFormatsValue = nil
		} else {
			supportedFormatsValue = common.ToPtr(strings.Split(supportedFormatsStr, ","))
		}

		dataRow := model.AttributeInfoRequest{
			Version:          row[headersMap["version"]],
			AttributeId:      row[headersMap["attr_id"]],
			AttributePath:    row[headersMap["attr_path"]],
			SupportedFormats: supportedFormatsValue,
			Field:            common.ToPtr(row[headersMap["field"]]),
			DataType:         common.ToPtr(row[headersMap["data_type"]]),
			Input:            common.ToPtr(row[headersMap["input"]]),
			Values:           common.ToPtr(row[headersMap["values"]]),
			Protocol:         row[headersMap["protocol"]],
			Examples:         common.ToPtr(row[headersMap["example"]]),
			KeySetName:       common.ToPtr(row[headersMap["key_set_name"]]),
			Description:      common.ToPtr(row[headersMap["description"]]),
			Executor:         row[headersMap["executor"]],
		}

		dtoList = append(dtoList, dataRow)
	}

	return dtoList, nil
}

func (a attributeService) UpsertAttributes(file multipart.File) (bool, *zkerrors.ZkError) {
	dtoList, zkError := readCSVAndReturnData(file)
	if zkError != nil {
		return false, nil
	}
	attributeDtoList := make(model.AttributeDtoList, 0)
	mapExecutorToDtoList := make(map[string][]model.AttributeInfoRequest)
	for _, v := range dtoList {
		if _, ok := mapExecutorToDtoList[v.Executor]; !ok {
			mapExecutorToDtoList[v.Executor] = make([]model.AttributeInfoRequest, 0)
		}
		mapExecutorToDtoList[v.Executor] = append(mapExecutorToDtoList[v.Executor], v)
	}

	for _, dtoList := range mapExecutorToDtoList {
		if valid, zkErr := validation.IsAttributesListValid(dtoList); !valid {
			return false, zkErr
		}

		l := model.ConvertAttributeInfoRequestToAttributeDto(dtoList)
		attributeDtoList = append(attributeDtoList, l...)
	}

	done, err := a.repo.UpsertAttributes(attributeDtoList)
	if err != nil {
		return false, nil
	}

	return done, nil
}
