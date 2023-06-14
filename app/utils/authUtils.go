package utils

import (
	"github.com/kataras/iris/v12"
	_ "github.com/zerok-ai/zk-utils-go/http"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"main/app/utils/errors"
)

func ValidateApiKeyMiddleware(ctx iris.Context) {
	request := ctx.Request()
	apiKey := request.Header["Zk-Api-Key"]

	if apiKey == nil || apiKey[0] == "" {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestZkApiKeyMiddlewareEmpty, nil)
		zkHttpResponse := zkHttp.ToZkResponse[any](zkErr.Error.Status, nil, nil, &zkErr)
		ctx.StopWithJSON(zkHttpResponse.Status, zkHttpResponse)
		return
	}

	ctx.Next()
}
