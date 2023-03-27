package cluster

import (
	"github.com/kataras/iris/v12"
	"main/app/cluster/models"
	"main/app/utils"
)

func ValidCluster(clusterId string) bool {
	// must be present in ClusterMap.
	_, exist := models.ClusterMap[clusterId]
	return exist
}

func ValidateResource(resource string, ctx iris.Context) bool {
	if !utils.Contains(utils.ResourceList, resource) {
		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid resource"))
		return false
	}
	return true
}

func ValidatePxlTime(ctx iris.Context, s string) bool {
	if !utils.IsValidPxlTime(s) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrPxlStartTimeEmpty)
		return false
	}
	return true
}

func ValidateCluster(clusterIdx string, ctx iris.Context) bool {
	if !ValidCluster(clusterIdx) {
		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid cluster ID"))
		return false
	}
	return true
}

func ValidateGraphDetailsApi(ctx iris.Context, serviceName, ns, st, apiKey string) bool {
	if utils.IsEmpty(serviceName) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrServiceNameEmpty)
		return false
	}
	if utils.IsEmpty(ns) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrNamespaceEmpty)
		return false
	}
	if utils.IsEmpty(st) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrPxlStartTimeEmpty)
		return false
	}
	if utils.IsEmpty(apiKey) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrZkApiKeyEmpty)
		return false
	}
	return true
}

func ValidateGetResourceDetailsApi(ctx iris.Context, st string, apiKey string) bool {
	if utils.IsEmpty(st) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrPxlStartTimeEmpty)
		return false
	}
	if utils.IsEmpty(apiKey) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrZkApiKeyEmpty)
		return false
	}
	return true
}

func ValidateGetPxlData(ctx iris.Context, s string, apiKey string) bool {
	if utils.IsEmpty(s) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrClusterIdEmpty)
		return false
	}
	if utils.IsEmpty(apiKey) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrZkApiKeyEmpty)
		return false
	}
	return true
}
