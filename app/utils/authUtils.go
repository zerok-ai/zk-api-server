package utils

import (
	"github.com/kataras/iris/v12"
	"main/app/utils/zkerrors"
)

func ValidateApiKeyMiddleware(ctx iris.Context) {
	var zkHttpResponse ZkHttpResponse
	request := ctx.Request()
	apiKey := request.Header["ZK_API_KEY"]

	if apiKey == nil {
		zkHttpResponse = ZkHttpResponseBuilder{}.WithZkErrorType(zkerrors.ZK_ERROR_SESSION_EXPIRED).
			Build()
		ctx.StopWithJSON(zkHttpResponse.Status, zkHttpResponse)
		//ctx.StatusCode(iris.StatusBadRequest)
		//ctx.SetErr(ErrZkApiKeyEmpty)
		return
	}

	ctx.Next()
}
