package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	attributeModel "zk-api-server/app/attribute/model"
	"zk-api-server/app/attribute/service"
	"zk-api-server/app/attribute/validation"
	"zk-api-server/app/integrations/model/transformer"
	"zk-api-server/app/utils"
	"zk-api-server/app/utils/errors"
	"zk-api-server/internal/model"
)

type AttributeHandler interface {
	GetAttributes(ctx iris.Context)
	GetAttributesForBackend(ctx iris.Context)
	UploadAttributesCSV(ctx iris.Context)
}

type attributeHandler struct {
	service service.AttributeService
	cfg     model.ZkApiServerConfig
}

func NewAttributeHandler(service service.AttributeService, cfg model.ZkApiServerConfig) AttributeHandler {
	return &attributeHandler{service, cfg}
}

func (a *attributeHandler) GetAttributes(ctx iris.Context) {
	protocol := ctx.URLParam(utils.Protocol)

	var zkHttpResponse zkHttp.ZkHttpResponse[attributeModel.AttributeListResponse]
	var zkErr *zkerrors.ZkError
	var resp *attributeModel.AttributeListResponse

	if zkErr = validation.ValidateGetAttributes(protocol); zkErr != nil {
		zkLogger.Error("Error while validating GetAttributes api, protocol is empty", protocol)
	} else {
		resp, zkErr = a.service.GetAttributes(protocol)
	}

	if a.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[attributeModel.AttributeListResponse](200, *resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[attributeModel.AttributeListResponse](200, *resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (a *attributeHandler) GetAttributesForBackend(ctx iris.Context) {
	version := ctx.URLParam(utils.Version)

	var zkHttpResponse zkHttp.ZkHttpResponse[attributeModel.ExecutorAttributesResponse]
	var zkErr *zkerrors.ZkError
	var resp *attributeModel.ExecutorAttributesResponse

	if zkErr = validation.ValidateGetAttributes(version); zkErr != nil {
		zkLogger.Error("Error while validating GetAttributes api, version is empty", version)
		zkErr = common.ToPtr(zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestVersionEmpty, nil))
	} else {
		resp, zkErr = a.service.GetAttributesForBackend(version)
	}

	if resp == nil {
		resp = &attributeModel.ExecutorAttributesResponse{}
	}

	if a.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[attributeModel.ExecutorAttributesResponse](200, *resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[attributeModel.ExecutorAttributesResponse](200, *resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (a *attributeHandler) UploadAttributesCSV(ctx iris.Context) {
	var zkHttpResponse zkHttp.ZkHttpResponse[any]
	var zkErr *zkerrors.ZkError
	done := false

	if file, _, err := ctx.Request().FormFile(utils.File); err != nil {
		zkLogger.Error("Error Retrieving the File", err)
		zkErr = common.ToPtr(zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestFileAttachedError, nil))
		zkHttpResponse = zkHttp.ZkHttpResponseBuilder[any]{}.WithZkErrorType(zkErr.Error).Build()
	} else {
		done, zkErr = a.service.UpsertAttributes(file)
	}

	if a.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[any](200, transformer.IntegrationResponse{}, done, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[any](200, transformer.IntegrationResponse{}, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}
