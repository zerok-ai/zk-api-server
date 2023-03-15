package px

import (
	"github.com/kataras/iris/v12"
)

func GetPxData(ctx iris.Context) {
	st := ctx.URLParamDefault("st", "-10m")
	clusterIdx := ctx.URLParam("cluster_id")
	if !ValidateGetPxlData(ctx, clusterIdx) {
		return
	}

	getPxlData(ctx, clusterIdx, st)
}
