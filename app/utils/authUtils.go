package utils

import (
	"github.com/kataras/iris/v12"
	"main/app/utils/zkerrors"
)

func ValidateApiKeyMiddleware(ctx iris.Context) {
	var zkHttpResponse ZkHttpResponse[any]
	request := ctx.Request()
	apiKey := request.Header["Zk_api_key"]

	if apiKey == nil || apiKey[0] == "" {
		z := ZkHttpResponseBuilder[any]{}
		zkHttpResponse = z.WithZkErrorType(zkerrors.ZK_ERROR_BAD_REQUEST_ZK_API_KEY_MIDDLEWARE_EMPTY).
			Build()
		ctx.StopWithJSON(zkHttpResponse.Status, zkHttpResponse)
		return
	}

	ctx.Next()
}
