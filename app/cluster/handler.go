package cluster

import (
	"github.com/kataras/iris/v12"
	"main/app/cluster/models"
	"main/app/utils"
)

func ListCluster(ctx iris.Context) {
	listCluster(ctx)
}

func UpsertCluster(ctx iris.Context) {
	var cluster models.Cluster
	err := ctx.ReadJSON(&cluster)
	// Validation: Is cluster model valid, parse cluster obj
	if err != nil {
		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Failed to parse cluster info").DetailErr(err))
		return
	}

	updateCluster(ctx, cluster)
}

func DeleteCluster(ctx iris.Context) {
	clusterId := ctx.Params().Get("clusterId")
	if utils.IsEmpty(clusterId) {
		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("clusterId cannot be empty"))
		return
	}

	deleteCluster(ctx, clusterId)
}

func GetResourceDetailsList(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	st := ctx.URLParam("st")
	ns := ctx.URLParam("ns")

	if !ValidateGetResourceDetailsApi(ctx, ns, st) {
		return
	}

	getResourceDetails(ctx, clusterIdx, ns, "list", st)
}

func GetResourceDetailsMap(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	st := ctx.URLParam("st")
	ns := ctx.URLParam("ns")

	if !ValidateGetResourceDetailsApi(ctx, ns, st) {
		return
	}

	getResourceDetails(ctx, clusterIdx, ns, "map", st)
}

func GetServiceStatsGraph(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	serviceName := ctx.URLParam("name")
	ns := ctx.URLParam("ns")
	st := ctx.URLParam("st")

	if !ValidateGraphStatsApi(ctx, serviceName, ns, st) {
		return
	}

	getServiceStatsGraph(ctx, clusterIdx, serviceName, ns, st)

}
