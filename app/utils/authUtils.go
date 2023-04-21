package utils

import (
	"github.com/kataras/iris/v12"
	"main/app/utils/zkerrors"
)

func ValidateApiKeyMiddleware(ctx iris.Context) {
	var zkHttpResponse ZkHttpResponse
	request := ctx.Request()
	apiKey := request.Header["Zk_api_key"]

	if apiKey == nil || apiKey[0] == "" {
		zkHttpResponse = ZkHttpResponseBuilder{}.WithZkErrorType(zkerrors.ZK_ERROR_BAD_REQUEST_ZK_API_KEY_EMPTY).
			Build()
		ctx.StopWithJSON(zkHttpResponse.Status, zkHttpResponse)
		return
	}

	ctx.Next()
}
