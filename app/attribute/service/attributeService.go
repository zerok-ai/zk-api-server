package service

import (
	"github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"log"
	"mime/multipart"
	"strconv"
	"zk-api-server/app/attribute/model"
	"zk-api-server/app/attribute/repository"
	"zk-api-server/utils"
)

type AttributeService interface {
	GetAttributes(protocol string) (*model.AttributeListResponse, *zkerrors.ZkError)
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

func (a attributeService) GetAttributes(protocol string) (*model.AttributeListResponse, *zkerrors.ZkError) {
	if common.IsEmpty(protocol) {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, nil)
		zkLogger.Error(LogTag, "protocol is empty")
		return nil, &zkError
	}

	data, err := a.repo.GetAttributes(protocol)
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
	if response.ExecutorAttributesList.Version > updatedAtInt {
		response.ExecutorAttributesList.Update = true
	}

	return &response, nil
}

func (a attributeService) UpsertAttributes(file multipart.File) (bool, *zkerrors.ZkError) {
	if file == nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, nil)
		zkLogger.Error(LogTag, "file is nil")
		return false, &zkError
	}

	csvData, err := utils.ParseCSV(file)
	dtoList := make([]model.AttributeInfoRequest, 0)
	for i, row := range csvData {
		sendToFrontEnd, _ := strconv.ParseBool(row[13])
		row := model.AttributeInfoRequest{
			Version:          row[0],
			CommonId:         row[1],
			VersionId:        row[2],
			Field:            row[3],
			DataType:         row[4],
			Input:            row[5],
			Values:           row[6],
			Protocol:         row[7],
			Examples:         row[8],
			KeySetName:       row[9],
			Description:      row[10],
			RequirementLevel: row[11],
			Executor:         row[12],
			SendToFrontEnd:   sendToFrontEnd,
		}

		if common.IsEmpty(row.VersionId) || common.IsEmpty(row.CommonId) || common.IsEmpty(row.Field) ||
			common.IsEmpty(row.DataType) || common.IsEmpty(row.Input) || common.IsEmpty(row.Protocol) {
			zkLogger.Error(LogTag, "missing required fields in csv file, line num: %d", i)
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, nil)
			return false, &zkError
		}

		dtoList = append(dtoList, row)
	}
	dtoListWithoutHeader := dtoList[1:]

	attributeDtoList := model.ConvertAttributeInfoRequestToAttributeDto(dtoListWithoutHeader)

	if err != nil {
		log.Println("Error parsing CSV:", err)
		return false, nil
	}

	done, err := a.repo.UpsertAttributes(attributeDtoList)
	if err != nil {
		return false, nil
	}

	return done, nil
}
