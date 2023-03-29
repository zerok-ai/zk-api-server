package cluster

import (
	"github.com/kataras/iris/v12"
	"main/app/cluster/models"
	"main/app/utils"
)

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
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParam("st")

	if !ValidateGetResourceDetailsApi(ctx, st, apiKey) {
		return
	}

	getResourceDetails(ctx, clusterIdx, "list", st, apiKey)
}

func GetResourceDetailsMap(ctx iris.Context) {
	apiKey := ctx.GetHeader("ZK_API_KEY")
	clusterIdx := ctx.Params().Get("clusterIdx")
	st := ctx.URLParam("st")

	if !ValidateGetResourceDetailsApi(ctx, st, apiKey) {
		return
	}

	getResourceDetails(ctx, clusterIdx, "map", st, apiKey)
}

func GetServiceDetails(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	apiKey := ctx.GetHeader("ZK_API_KEY")
	serviceName := ctx.URLParam("name")
	ns := ctx.URLParam("ns")
	st := ctx.URLParam("st")

	if !ValidateGraphDetailsApi(ctx, serviceName, ns, st, apiKey) {
		return
	}

	getServiceDetails(ctx, clusterIdx, serviceName, ns, st, apiKey)

}

func GetPodDetailsList(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParam("st")
	serviceName := ctx.URLParam("service_name")
	ns := ctx.URLParam("ns")

	if !ValidateGraphDetailsApi(ctx, serviceName, ns, st, apiKey) {
		return
	}

	getPodDetails(ctx, clusterIdx, serviceName, ns, st, apiKey)
}

func GetPodDetails(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParam("st")
	podName := ctx.URLParam("pod_name")
	ns := ctx.URLParam("ns")

	if !ValidatePodDetailsApi(ctx, podName, ns, st, apiKey) {
		return
	}

	getPodDetailsTimeSeries(ctx, clusterIdx, podName, ns, st, apiKey)
}

func GetPxData(ctx iris.Context) {
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParamDefault("st", "-10m")
	clusterIdx := ctx.URLParam("cluster_id")
	if !ValidateGetPxlData(ctx, clusterIdx, apiKey) {
		return
	}

	getPxlData(ctx, clusterIdx, st, apiKey)
}
