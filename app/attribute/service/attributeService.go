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

func (a attributeService) UpsertAttributes(file multipart.File) (bool, *zkerrors.ZkError) {
	if file == nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, nil)
		zkLogger.Error(LogTag, "file is nil")
		return false, &zkError
	}

	csvData, err := utils.ParseCSV(file)
	dtoList := make([]model.AttributeInfoRequest, 0)
	for _, row := range csvData {
		sendToFrontEnd, _ := strconv.ParseBool(row[12])
		dtoList = append(dtoList, model.AttributeInfoRequest{
			Version:          row[0],
			Id:               row[1],
			Field:            row[2],
			DataType:         row[3],
			Input:            row[4],
			Values:           row[5],
			Protocol:         row[6],
			Examples:         row[7],
			KeySetName:       row[8],
			Description:      row[9],
			RequirementLevel: row[10],
			Executor:         row[11],
			SendToFrontEnd:   sendToFrontEnd,
		})
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
