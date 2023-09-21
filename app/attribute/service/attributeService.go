package service

import (
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"log"
	"mime/multipart"
	"zk-api-server/app/attribute/model"
	"zk-api-server/app/attribute/repository"
	"zk-api-server/utils"
)

type AttributeService interface {
	GetAttributes(version string, keySet []string) (*model.AttributeListResponse, *zkerrors.ZkError)
	UpsertAttributes(multipart.File) (bool, *zkerrors.ZkError)
}

var LogTag = "attribute_service"

type attributeService struct {
	repo repository.AttributeRepo
}

func NewAttributeService(repo repository.AttributeRepo) AttributeService {
	return &attributeService{repo: repo}
}

func (a attributeService) GetAttributes(version string, keySet []string) (*model.AttributeListResponse, *zkerrors.ZkError) {
	if version == "" {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, nil)
		zkLogger.Error(LogTag, "version is empty")
		return nil, &zkError
	}

	data, err := a.repo.GetAttributes(version, keySet)
	if err != nil {
		zkLogger.Error(LogTag, "failed to get attributes list", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorDbError, err)
		return nil, &zkError

	}
	response := model.ConvertAttributeDtoToAttributeResponse(data, version)
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
		dtoList = append(dtoList, model.AttributeInfoRequest{
			Version:          row[0],
			Attribute:        row[1],
			KeySetName:       row[2],
			Type:             row[3],
			Description:      row[4],
			Examples:         row[5],
			RequirementLevel: row[6],
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
