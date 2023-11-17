package service

import (
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	dto "zk-api-server/app/obfuscation/model/dto"
	"zk-api-server/app/obfuscation/model/transformer"
	"zk-api-server/app/obfuscation/repository"
)

type ObfuscationService interface {
	GetAllObfuscationsDashboard(orgId string, offset, limit string) (transformer.ObfuscationListResponse, *zkerrors.ZkError)
	GetAllObfuscationsOperator(orgId string, updatedTime int64) (transformer.ObfuscationResponseOperator, *zkerrors.ZkError)
	GetObfuscationById(id string, orgId string) (transformer.ObfuscationResponse, *zkerrors.ZkError)
	InsertObfuscation(obfuscation dto.Obfuscation) (bool, *zkerrors.ZkError)
	UpdateObfuscation(obfuscation dto.Obfuscation) (bool, *zkerrors.ZkError)
	DeleteObfuscation(orgId string, id string) (bool, *zkerrors.ZkError)
}

var LogTag = "obfuscation_service"

type obfuscationService struct {
	repo repository.ObfuscationRepo
}

func (o obfuscationService) GetAllObfuscationsOperator(orgId string, updatedTime int64) (transformer.ObfuscationResponseOperator, *zkerrors.ZkError) {
	obfuscations, err := o.repo.GetAllObfuscationsOperator(orgId, updatedTime)
	if err != nil {
		zkLogger.Error(LogTag, "Error while getting all obfuscations: ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return transformer.ObfuscationResponseOperator{}, &zkError
	}
	return transformer.ToObfuscationListResponseOperator(obfuscations), nil
}

func NewObfuscationService(repo repository.ObfuscationRepo) ObfuscationService {
	return &obfuscationService{repo: repo}
}

func (o obfuscationService) GetAllObfuscationsDashboard(orgId string, offset, limit string) (transformer.ObfuscationListResponse, *zkerrors.ZkError) {

	obfuscations, totalRows, err := o.repo.GetAllObfuscationsForDashboard(orgId, offset, limit)
	if err != nil {
		zkLogger.Error(LogTag, "Error while getting all obfuscations: ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return transformer.ObfuscationListResponse{}, &zkError
	}

	return transformer.ToObfuscationListResponse(obfuscations, totalRows), nil
}

func (o obfuscationService) GetObfuscationById(id string, orgId string) (transformer.ObfuscationResponse, *zkerrors.ZkError) {
	obfuscation, err := o.repo.GetObfuscationById(id, orgId)
	if err != nil {
		zkLogger.Error(LogTag, "Error while getting obfuscation by ID: ", id, err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return transformer.ObfuscationResponse{Response: nil}, &zkError
	}

	if obfuscation == nil || obfuscation.Deleted {
		zkLogger.Error(LogTag, "Getting obfuscation nil for by ID: ", id, err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, "Obfuscation rule not found for the given id.")
		return transformer.ObfuscationResponse{Response: nil}, &zkError
	}

	response, err := transformer.ToObfuscationResponse(*obfuscation)

	if err != nil {
		zkLogger.Error(LogTag, "Error while converting obfuscation to response: ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return transformer.ObfuscationResponse{Response: nil}, &zkError
	}

	return response, nil
}

func (o obfuscationService) InsertObfuscation(obfuscation dto.Obfuscation) (bool, *zkerrors.ZkError) {
	done, err := o.repo.InsertObfuscation(obfuscation)
	if err != nil {
		zkLogger.Error(LogTag, "Error while inserting obfuscation: ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return false, &zkError
	}

	return done, nil
}

func (o obfuscationService) UpdateObfuscation(obfuscation dto.Obfuscation) (bool, *zkerrors.ZkError) {
	done, err := o.repo.UpdateObfuscation(obfuscation)
	if err != nil {
		zkLogger.Error(LogTag, "Error while updating obfuscation: ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return false, &zkError
	}

	return done, nil
}

func (o obfuscationService) DeleteObfuscation(orgId string, id string) (bool, *zkerrors.ZkError) {
	done, err := o.repo.DeleteObfuscation(orgId, id)
	if err != nil {
		zkLogger.Error(LogTag, "Error while deleting obfuscation: ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return false, &zkError
	}

	if !done {
		zkLogger.Error(LogTag, "Failed to delete obfuscation.")
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return false, &zkError
	}

	return true, nil
}
