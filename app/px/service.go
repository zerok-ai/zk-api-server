package px

import (
	"github.com/kataras/iris/v12"
	"main/app/tablemux"
	"main/app/tablemux/handlerimplementation"
	"main/app/utils"
)

func getPxlData(ctx iris.Context, clusterIdx, st string) {

	var s = make([]handlerimplementation.PixieTraceData, 0)
	pixieTraceDataMux := handlerimplementation.PixieTraceDataListMux{Table: handlerimplementation.TablePrinterPixieTraceDataList{Values: s}}

	tx := tablemux.MethodTemplate{MethodSignature: utils.GetPXDataSignature(100, st, "{}"), DataFrameName: "my_first_list"}
	resultSet := tablemux.GetResource(ctx, clusterIdx, &pixieTraceDataMux, tx, 3)
	result := pixieTraceDataMux.Table.Values

	if result == nil {
		return
	}

	ctx.StatusCode(iris.StatusOK)
	_ = ctx.JSON(map[string]interface{}{
		"results": result,
		"stats":   resultSet.Stats(),
		"status":  200,
	})
}
