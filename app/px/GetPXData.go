package px

import (
	"github.com/kataras/iris/v12"
	"main/app/cluster"
	"main/app/cluster/models"
	"main/app/cluster/models/tablemux"
	"main/app/utils"
	"px.dev/pxapi"
)

type Template struct {
	StartTime string
	Head      int
	Filter    string
}

func getPXData(ctx iris.Context, clusterDetails models.Cluster) (*pxapi.ScriptResults, []models.PixieTraceData) {
	var s = make([]models.PixieTraceData, 0)
	pixieTraceDataMux := tablemux.PixieTraceDataListMux{Table: tablemux.TablePrinterPixieTraceDataList{Values: s}}
	//pixieTraceDataMux := tableMux{}

	tx := models.MethodTemplate{MethodSignature: utils.GetPXDataSignature(100, "-10m", "{}"), DataFrameName: "my_first_list"}
	resultSet := cluster.GetResource(ctx, clusterDetails, &pixieTraceDataMux, tx, 3)
	return resultSet, pixieTraceDataMux.Table.Values
}

func GetPXData(ctx iris.Context) {
	clusterMapId := ctx.URLParamDefault("cluster_id", "1")
	clusterDetails := cluster.GetClusterDetails(clusterMapId)

	var resultSet *pxapi.ScriptResults
	var result interface{}

	resultSet, result = getPXData(ctx, clusterDetails)

	if result == nil {
		return
	}

	_ = ctx.JSON(map[string]interface{}{
		"results": result,
		"stats":   resultSet.Stats(),
		"status":  200,
	})
}
