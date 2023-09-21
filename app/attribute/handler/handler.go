package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strings"
	model2 "zk-api-server/app/attribute/model"
	"zk-api-server/app/attribute/service"
	"zk-api-server/app/attribute/validation"
	"zk-api-server/app/integrations/model/transformer"
	"zk-api-server/app/utils/errors"
	"zk-api-server/internal/model"
)

type AttributeHandler interface {
	GetAttributes(ctx iris.Context)
	UploadCSVHandler(ctx iris.Context)
}

type attributeHandler struct {
	service service.AttributeService
	cfg     model.ZkApiServerConfig
}

func NewAttributeHandler(service service.AttributeService, cfg model.ZkApiServerConfig) AttributeHandler {
	return &attributeHandler{service, cfg}
}

func (a *attributeHandler) GetAttributes(ctx iris.Context) {
	version := ctx.URLParam("version")
	keySets := ctx.URLParam("key_set")

	var zkHttpResponse zkHttp.ZkHttpResponse[model2.AttributeListResponse]
	var zkErr *zkerrors.ZkError
	var resp *model2.AttributeListResponse

	if zkErr = validation.ValidateGetAttributes(version, keySets); zkErr != nil {
		zkLogger.Error("Error while validating GetAttributes api, version: %s or keySets: %s is empty", version, keySets)
	} else {
		if common.IsEmpty(keySets) {
			resp, zkErr = a.service.GetAttributes(version, nil)
		} else {
			k := strings.Split(keySets, ",")
			resp, zkErr = a.service.GetAttributes(version, k)
		}
	}

	if a.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[model2.AttributeListResponse](200, *resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[model2.AttributeListResponse](200, *resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (a *attributeHandler) UploadCSVHandler(ctx iris.Context) {
	var zkHttpResponse zkHttp.ZkHttpResponse[any]
	var zkErr *zkerrors.ZkError
	done := false

	if file, _, err := ctx.Request().FormFile("file"); err != nil {
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
