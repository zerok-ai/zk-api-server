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

//func ValidateResource(resource string, ctx iris.Context) bool {
//	if !utils.Contains(utils.ResourceList, resource) {
//		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
//			Title("Invalid resource"))
//		return false
//	}
//	return true
//}

func ValidateAction(action string, ctx iris.Context) bool {
	if !utils.Contains(utils.Actions, action) {
		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid action"))
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

func ValidateGraphStatsApi(ctx iris.Context, serviceName, ns, st string) bool {
	if utils.IsEmpty(serviceName) {
		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("service name cannot be empty"))
		return false
	}
	if utils.IsEmpty(ns) {
		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("namespace cannot be empty"))
		return false
	}
	if utils.IsEmpty(st) {
		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("start time cannot be empty"))
		return false
	}
	return true
}

func ValidateGetResourceDetailsApi(ctx iris.Context, ns, st string) bool {
	if utils.IsEmpty(ns) {
		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("namespace cannot be empty"))
		return false
	}
	if utils.IsEmpty(st) {
		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("start time cannot be empty"))
		return false
	}
	return true
}
