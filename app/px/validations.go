package px

import (
	"github.com/kataras/iris/v12"
	"main/app/utils"
)

func ValidateGetPxlData(ctx iris.Context, s string) bool {
	if utils.IsEmpty(s) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrClusterIdEmpty)
		return false
	}
	return true
}
