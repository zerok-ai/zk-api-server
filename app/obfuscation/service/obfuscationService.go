package service

import (
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	dto "zk-api-server/app/obfuscation/model/dto"
	"zk-api-server/app/obfuscation/model/transformer"
	"zk-api-server/app/obfuscation/repository"
)

type ObfuscationService interface {
	GetAllObfuscations(orgId string, offset, limit string) (transformer.ObfuscationListResponse, *zkerrors.ZkError)
	GetObfuscationById(id string, orgId string) (transformer.ObfuscationResponse, *zkerrors.ZkError)
	InsertObfuscation(obfuscation dto.Obfuscation) (bool, *zkerrors.ZkError)
	UpdateObfuscation(obfuscation dto.Obfuscation) (bool, *zkerrors.ZkError)
	DeleteObfuscation(orgId string, id string) (bool, *zkerrors.ZkError)
}

var LogTag = "obfuscation_service"

type obfuscationService struct {
	repo repository.ObfuscationRepo
}

func NewObfuscationService(repo repository.ObfuscationRepo) ObfuscationService {
	return &obfuscationService{repo: repo}
}

func (o obfuscationService) GetAllObfuscations(orgId string, offset, limit string) (transformer.ObfuscationListResponse, *zkerrors.ZkError) {
	obfuscations, err := o.repo.GetAllObfuscations(orgId, offset, limit)
	if err != nil {
		zkLogger.Error(LogTag, "Error while getting all obfuscations: ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return transformer.ObfuscationListResponse{}, &zkError
	}

	return transformer.ToObfuscationListResponse(obfuscations), nil
}

func (o obfuscationService) GetObfuscationById(id string, orgId string) (transformer.ObfuscationResponse, *zkerrors.ZkError) {
	obfuscation, err := o.repo.GetObfuscationById(id, orgId)
	if err != nil {
		zkLogger.Error(LogTag, "Error while getting obfuscation by ID: ", id, err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return transformer.ObfuscationResponse{}, &zkError
	}

	if obfuscation == nil {
		zkLogger.Error(LogTag, "Getting obfuscation nil for by ID: ", id, err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return transformer.ObfuscationResponse{}, &zkError
	}

	response, err := transformer.ToObfuscationResponse(*obfuscation)

	if err != nil {
		zkLogger.Error(LogTag, "Error while converting obfuscation to response: ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return transformer.ObfuscationResponse{}, &zkError
	}

	return response, nil
}

func (o obfuscationService) InsertObfuscation(obfuscation dto.Obfuscation) (bool, *zkerrors.ZkError) {
	//TODO: Perform validation on regex here.
	done, err := o.repo.InsertObfuscation(obfuscation)
	if err != nil {
		zkLogger.Error(LogTag, "Error while inserting obfuscation: ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return false, &zkError
	}

	return done, nil
}

func (o obfuscationService) UpdateObfuscation(obfuscation dto.Obfuscation) (bool, *zkerrors.ZkError) {
	//TODO: Perform validation on regex here.
	done, err := o.repo.UpdateObfuscation(obfuscation)
	if err != nil {
		zkLogger.Error(LogTag, "Error while updating obfuscation: ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return false, &zkError
	}

	return done, nil
}

func (o obfuscationService) DeleteObfuscation(orgId string, id string) (bool, *zkerrors.ZkError) {
	done, err := o.repo.DeleteObfuscation(orgId, id)
	if err != nil {
		zkLogger.Error(LogTag, "Error while deleting obfuscation: ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return false, &zkError
	}

	if !done {
		zkLogger.Error(LogTag, "Failed to delete obfuscation.")
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return false, &zkError
	}

	return true, nil
}
