package utils

import (
	"github.com/kataras/iris/v12"
	"main/app/cluster/models"
	"main/app/utils/zkerrors"
)

// TODO: Move to zkhttp
type ZkHttpError struct {
	Kind    string          `json:"kind,omitempty"`
	Message string          `json:"message,omitempty"`
	Param   string          `json:"param,omitempty"`
	Stack   any             `json:"stack,omitempty"`
	Info    *map[string]any `json:"info,omitempty"`
}

type ZkHttpResponse struct {
	Metadata *map[string]any    `json:"-"`
	Headers  *map[string]string `json:"-"`
	Status   int                `json:"-"`
	Error    *ZkHttpError       `json:"error,omitempty"`
	Message  *string            `json:"message,omitempty"`
	Debug    *map[string]any    `json:"debug,omitempty"`
	Data     any                `json:"payload,omitempty"`
}

// Zk Http Response Builder
type ZkHttpResponseBuilder struct {
	_ZkHttpResponseBuilder _zkHttpResponseBuilder
}

func (zkHttpResponseBuilder ZkHttpResponseBuilder) WithStatus(status int) _zkHttpResponseBuilder {
	zkHttpResponseBuilder._ZkHttpResponseBuilder = _zkHttpResponseBuilder{}
	return zkHttpResponseBuilder._ZkHttpResponseBuilder.withStatus(status)
}

func (zkHttpResponseBuilder ZkHttpResponseBuilder) WithZkErrorType(zkErrorType zkerrors.ZkErrorType) _zkHttpResponseBuilder {
	zkHttpResponseBuilder._ZkHttpResponseBuilder = _zkHttpResponseBuilder{}
	zkHttpError := ZkHttpError{}
	zkHttpError.Build(zkErrorType, nil, nil)
	return zkHttpResponseBuilder._ZkHttpResponseBuilder.withStatus(zkErrorType.Status).Error(zkHttpError)
}

// Zk Http Response Builder Internal
type _zkHttpResponseBuilder struct {
	ZkHttpResponse ZkHttpResponse
}

func (zkHttpResponseBuilder _zkHttpResponseBuilder) withStatus(status int) _zkHttpResponseBuilder {
	zkHttpResponseBuilder.ZkHttpResponse.Status = status
	return zkHttpResponseBuilder
}

func (zkHttpResponseBuilder _zkHttpResponseBuilder) Message(message *string) _zkHttpResponseBuilder {
	zkHttpResponseBuilder.ZkHttpResponse.Message = message
	return zkHttpResponseBuilder
}

func (zkHttpResponseBuilder _zkHttpResponseBuilder) Data(data any) _zkHttpResponseBuilder {
	if data != nil {
		zkHttpResponseBuilder.ZkHttpResponse.Data = data
	}
	return zkHttpResponseBuilder
}

func (zkHttpResponseBuilder _zkHttpResponseBuilder) Debug(key string, value any) _zkHttpResponseBuilder {
	if !HTTP_DEBUG {
		return zkHttpResponseBuilder
	}
	if zkHttpResponseBuilder.ZkHttpResponse.Debug == nil {
		zkHttpResponseBuilder.ZkHttpResponse.Debug = &map[string]any{}
	}
	(*zkHttpResponseBuilder.ZkHttpResponse.Debug)[key] = value
	return zkHttpResponseBuilder
}

// func (zkHttpResponseBuilder *_zkHttpResponseBuilder) DebugP(key string, value any) *_zkHttpResponseBuilder{
// 	if !utils.HTTP_DEBUG {
// 		return zkHttpResponseBuilder;
// 	}
// 	if zkHttpResponseBuilder.ZkHttpResponse.Debug == nil {
// 		zkHttpResponseBuilder.ZkHttpResponse.Debug = &map[string]any{}
// 	}
// 	(*zkHttpResponseBuilder.ZkHttpResponse.Debug)[key] = value
// 	return zkHttpResponseBuilder;
// }

func (zkHttpResponseBuilder _zkHttpResponseBuilder) Header(key string, value string) _zkHttpResponseBuilder {
	if zkHttpResponseBuilder.ZkHttpResponse.Headers == nil {
		zkHttpResponseBuilder.ZkHttpResponse.Headers = &map[string]string{}
	}
	(*zkHttpResponseBuilder.ZkHttpResponse.Headers)[key] = value
	return zkHttpResponseBuilder
}

func (zkHttpResponseBuilder _zkHttpResponseBuilder) Metadata(key string, value any) _zkHttpResponseBuilder {
	if zkHttpResponseBuilder.ZkHttpResponse.Metadata == nil {
		zkHttpResponseBuilder.ZkHttpResponse.Metadata = &map[string]any{}
	}
	(*zkHttpResponseBuilder.ZkHttpResponse.Metadata)[key] = value
	return zkHttpResponseBuilder
}

func (zkHttpResponseBuilder _zkHttpResponseBuilder) Error(zkHttpError ZkHttpError) _zkHttpResponseBuilder {
	if zkHttpResponseBuilder.ZkHttpResponse.Error == nil {
		zkHttpResponseBuilder.ZkHttpResponse.Error = &zkHttpError
	}
	return zkHttpResponseBuilder
}

func (zkHttpResponseBuilder _zkHttpResponseBuilder) ErrorInfo(key string, value any) _zkHttpResponseBuilder {
	if zkHttpResponseBuilder.ZkHttpResponse.Error == nil {
		zkHttpResponseBuilder.ZkHttpResponse.Error = &ZkHttpError{}
	}
	if zkHttpResponseBuilder.ZkHttpResponse.Error.Info == nil {
		zkHttpResponseBuilder.ZkHttpResponse.Error.Info = &map[string]any{}
	}
	(*zkHttpResponseBuilder.ZkHttpResponse.Error.Info)[key] = value
	return zkHttpResponseBuilder
}

func (zkHttpResponseBuilder _zkHttpResponseBuilder) Build() ZkHttpResponse {
	return zkHttpResponseBuilder.ZkHttpResponse
}

func (zkHttpError *ZkHttpError) Build(zkErrorType zkerrors.ZkErrorType, param *string, stack any) {
	zkHttpError.Kind = zkErrorType.Type
	zkHttpError.Message = zkErrorType.Message
	if HTTP_DEBUG {
		zkHttpError.Stack = stack
	}
	if param != nil {
		zkHttpError.Param = *param
	}
}

func (zkHttpResponse ZkHttpResponse) Header(key string, value string) {
	if zkHttpResponse.Headers == nil {
		zkHttpResponse.Headers = &map[string]string{}
	}
	(*zkHttpResponse.Headers)[key] = value
}

func (zkHttpResponse ZkHttpResponse) IsOk() bool {
	if zkHttpResponse.Status > 199 && zkHttpResponse.Status < 300 {
		return true
	}
	return false
}
func GenerateResponseAndReturn(ctx iris.Context, pxResp models.PixieResponse) {

	var zkHttpResponse ZkHttpResponse
	if pxResp.Error != nil {
		zkHttpResponse = CreateErrorResponseWithStatusCode(pxResp.Error.Error)
	} else {
		zkHttpResponse = CreateSuccessResponseWithStatusCode(pxResp, 200)
	}
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}
func CreateSuccessResponseWithStatusCode(resp interface{}, statusCode int) ZkHttpResponse {
	zkHttpResponseBuilder := ZkHttpResponseBuilder{}
	var zkHttpResponse ZkHttpResponse
	zkHttpResponse = zkHttpResponseBuilder.WithStatus(statusCode).
		Message(nil).
		Data(resp).
		Build()
	return zkHttpResponse
}

func CreateErrorResponseWithStatusCode(zkErrorType zkerrors.ZkErrorType) ZkHttpResponse {
	var zkHttpResponse ZkHttpResponse
	zkHttpResponse = ZkHttpResponseBuilder{}.WithZkErrorType(zkErrorType).
		Build()
	return zkHttpResponse
}
