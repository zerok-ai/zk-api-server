package handler

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkObfuscation "github.com/zerok-ai/zk-utils-go/obfuscation/model"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
	"zk-api-server/app/obfuscation/model/transformer"
	"zk-api-server/app/obfuscation/service"
	"zk-api-server/app/obfuscation/validation"
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
	GetAllRulesOperator(ctx iris.Context)
}

type obfuscationHandler struct {
	service service.ObfuscationService
	cfg     zkApiModel.ZkApiServerConfig
}

func (o obfuscationHandler) GetAllRulesOperator(ctx iris.Context) {
	orgId := ctx.GetHeader(HTTP_UTILS_ORG_ID)
	if zkCommon.IsEmpty(orgId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("OrgId is required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[transformer.ObfuscationResponseOperator]
	var zkErr *zkerrors.ZkError
	var resp transformer.ObfuscationResponseOperator

	updateTimeStr := ctx.URLParamDefault(utils.LastSyncTS, "0")
	updateTime := parseIntWithDefault(updateTimeStr, 0)

	resp, zkErr = o.service.GetAllObfuscationsOperator(orgId, updateTime) // You can specify offset and limit as needed

	if o.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ObfuscationResponseOperator](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ObfuscationResponseOperator](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func NewObfuscationHandler(s service.ObfuscationService, cfg zkApiModel.ZkApiServerConfig) ObfuscationHandler {
	return &obfuscationHandler{service: s, cfg: cfg}
}

// GetAllRulesDashboard godoc
// @Summary Get all non-deleted obfuscation rules for an organization
// @Description Get all non-deleted obfuscation rules for an organization
// @Tags Obfuscation
// @Accept json
// @Produce json
// @Param orgId header string true "Organization ID"
// @Param limit query string false "Limit". Default 20
// @Param offset query string false "Offset". Default 0
// @Success 200 {object} ObfuscationListResponse
// @Failure 400 {string} string "OrgId is required"
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

	resp, zkErr = o.service.GetAllObfuscationsDashboard(orgId, offset, limit) // You can specify offset and limit as needed

	if o.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ObfuscationListResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[transformer.ObfuscationListResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

// GetObfuscationById godoc
// @Summary Get obfuscation rule by id
// @Description Get obfuscation rule by id.
// @Tags Obfuscation
// @Accept json
// @Produce json
// @Param orgId header string true "Organization ID"

func (o obfuscationHandler) GetObfuscationById(ctx iris.Context) {
	orgId := ctx.GetHeader(HTTP_UTILS_ORG_ID)
	id := ctx.Params().Get(utils.ObfuscationIdxPathParam)
	if zkCommon.IsEmpty(orgId) || zkCommon.IsEmpty(id) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("OrgId and Id are required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[transformer.ObfuscationResponse]

	resp, zkErr := o.service.GetObfuscationById(id, orgId)

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
	var zkError *zkerrors.ZkError

	done, message := validation.ValidateObfuscationRule(request)
	if !done {
		zkError2 := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, message)
		zkError = &zkError2
	} else {
		result, zkError2 := o.handlerUpsertInternal(ctx, isInsert, request, orgId)
		done = result
		zkError = zkError2
	}

	if o.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[bool](200, done, done, zkError)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[bool](200, done, nil, zkError)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (o obfuscationHandler) handlerUpsertInternal(ctx iris.Context, isInsert bool, request zkObfuscation.Rule, orgId string) (bool, *zkerrors.ZkError) {
	var done bool
	var zkError *zkerrors.ZkError
	if isInsert {
		done, zkError = o.service.InsertObfuscation(*transformer.FromObfuscationRequestToObfuscationDto(request, orgId, ""))
	} else {
		id := ctx.Params().Get(utils.ObfuscationIdxPathParam)
		if zkCommon.IsEmpty(orgId) {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString("Id is required")
			return false, nil
		}
		done, zkError = o.service.UpdateObfuscation(*transformer.FromObfuscationRequestToObfuscationDto(request, orgId, id))
	}
	return done, zkError
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

func parseIntWithDefault(s string, defaultVal int64) int64 {
	if s == "" {
		return defaultVal
	}

	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		// If there is an error during conversion, return the default value
		return defaultVal
	}
	return val
}
