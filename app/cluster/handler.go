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
	var cluster models.ClusterDetails
	err := ctx.ReadJSON(&cluster)
	// Validation: Is cluster model valid, parse cluster obj
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrClusterParsingFailed)
		return
	}

	updateCluster(ctx, cluster)
}

func DeleteCluster(ctx iris.Context) {
	clusterId := ctx.Params().Get("clusterId")

	if utils.IsEmpty(clusterId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrClusterIdEmpty)
		return
	}

	deleteCluster(ctx, clusterId)
}

func GetResourceDetailsList(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	st := ctx.URLParam("st")

	if !ValidateGetResourceDetailsApi(ctx, st) {
		return
	}

	getResourceDetails(ctx, clusterIdx, "list", st)
}

func GetResourceDetailsMap(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	st := ctx.URLParam("st")

	if !ValidateGetResourceDetailsApi(ctx, st) {
		return
	}

	getResourceDetails(ctx, clusterIdx, "map", st)
}

func GetServiceDetails(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	serviceName := ctx.URLParam("name")
	ns := ctx.URLParam("ns")
	st := ctx.URLParam("st")

	if !ValidateGraphDetailsApi(ctx, serviceName, ns, st) {
		return
	}

	getServiceDetails(ctx, clusterIdx, serviceName, ns, st)

}

func GetPodDetailsList(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	st := ctx.URLParam("st")
	serviceName := ctx.URLParam("service_name")
	ns := ctx.URLParam("ns")

	if !ValidateGraphDetailsApi(ctx, serviceName, ns, st) {
		return
	}

	getPodDetails(ctx, clusterIdx, serviceName, ns, st)
}
