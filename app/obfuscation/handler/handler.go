package handler

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkObfuscation "github.com/zerok-ai/zk-utils-go/obfuscation/model"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"zk-api-server/app/obfuscation/model/transformer"
	"zk-api-server/app/obfuscation/service"
	"zk-api-server/app/utils"
	zkApiModel "zk-api-server/internal/model"
)

var LogTag = "obfuscation_handler"
var HTTP_UTILS_ORG_ID = "X-ORG-ID"

type ObfuscationHandler interface {
	GetAllRulesDashboard(ctx iris.Context)
	GetObfuscationById(ctx iris.Context)
	InsertObfuscationRule(ctx iris.Context)
	UpdateObfuscationRule(ctx iris.Context)
	DeleteObfuscationRule(ctx iris.Context)
}

type obfuscationHandler struct {
	service service.ObfuscationService
	cfg     zkApiModel.ZkApiServerConfig
}

func NewObfuscationHandler(s service.ObfuscationService, cfg zkApiModel.ZkApiServerConfig) ObfuscationHandler {
	return &obfuscationHandler{service: s, cfg: cfg}
}

func (o obfuscationHandler) GetAllRulesDashboard(ctx iris.Context) {
	orgId := ctx.GetHeader(HTTP_UTILS_ORG_ID)
	if zkCommon.IsEmpty(orgId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("OrgId is required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[transformer.ObfuscationListResponse]
	var zkErr *zkerrors.ZkError
	var resp transformer.ObfuscationListResponse

	limit := ctx.URLParamDefault(utils.Limit, "20")
	offset := ctx.URLParamDefault(utils.Offset, "0")

	resp, zkErr = o.service.GetAllObfuscations(orgId, offset, limit) // You can specify offset and limit as needed

	if o.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ObfuscationListResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ObfuscationListResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (o obfuscationHandler) GetObfuscationById(ctx iris.Context) {
	orgId := ctx.GetHeader(HTTP_UTILS_ORG_ID)
	id := ctx.Params().Get(utils.ObfuscationIdxPathParam)
	if zkCommon.IsEmpty(orgId) || zkCommon.IsEmpty(id) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("OrgId and Id are required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[transformer.ObfuscationResponse]
	var zkErr *zkerrors.ZkError
	var resp transformer.ObfuscationResponse

	resp, zkErr = o.service.GetObfuscationById(id, orgId)

	if o.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ObfuscationResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ObfuscationResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (o obfuscationHandler) InsertObfuscationRule(ctx iris.Context) {
	o.UpsertObfuscation(ctx, true)
}

func (o obfuscationHandler) UpdateObfuscationRule(ctx iris.Context) {
	o.UpsertObfuscation(ctx, false)
}

func (o obfuscationHandler) UpsertObfuscation(ctx iris.Context, isInsert bool) {
	orgId := ctx.GetHeader(HTTP_UTILS_ORG_ID)
	if zkCommon.IsEmpty(orgId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("OrgId is required")
		return
	}

	var request zkObfuscation.Rule
	var zkHttpResponse zkHttp.ZkHttpResponse[bool]

	body, err := ctx.GetBody()
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Error reading request body")
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Error decoding JSON")
		return
	}

	//TODO: Add validation for request here. Especially for regex.

	var done bool
	var zkError *zkerrors.ZkError

	if isInsert {
		done, zkError = o.service.InsertObfuscation(*transformer.FromObfuscationRequestToObfuscationDto(request, orgId, ""))
	} else {
		id := ctx.Params().Get(utils.ObfuscationIdxPathParam)
		if zkCommon.IsEmpty(orgId) {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString("Id is required")
			return
		}
		done, zkError = o.service.UpdateObfuscation(*transformer.FromObfuscationRequestToObfuscationDto(request, orgId, id))
	}

	if o.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[bool](200, done, done, zkError)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[bool](200, done, nil, zkError)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (o obfuscationHandler) DeleteObfuscationRule(ctx iris.Context) {
	orgId := ctx.GetHeader(HTTP_UTILS_ORG_ID)
	id := ctx.Params().Get(utils.ObfuscationIdxPathParam)
	if zkCommon.IsEmpty(orgId) || zkCommon.IsEmpty(id) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("OrgId and Id are required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[bool]
	var zkErr *zkerrors.ZkError

	done, zkErr := o.service.DeleteObfuscation(orgId, id)

	if o.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[bool](200, done, done, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[bool](200, done, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}
