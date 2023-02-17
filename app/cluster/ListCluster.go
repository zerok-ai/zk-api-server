package cluster

import (
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func ListCluster(ctx iris.Context) {
	var clusters []Cluster
	for k, v := range ClusterMap {
		var id = k
		v.Id = &id
		clusters = append(clusters, v)
	}
	if clusters == nil {
		emptyArr := []string{}
		ctx.JSON(emptyArr)
		return
	}
	ctx.JSON(clusters)
}

func UpsertCluster(ctx iris.Context) {
	var cluster Cluster
	err := ctx.ReadJSON(&cluster)
	// Validation: Is cluster model valid, parse cluster obj
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Failed to parse cluster info").DetailErr(err))
		return
	}

	// Create a cluster if cluster.ID is missing. else, update if valid.
	if cluster.Id != nil {
		// Validation: Check if provided cluster ID is valid?
		if !ValidClusterId(*cluster.Id) {
			ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
				Title("Cluster with provided ID does not exist"))
			return
		}
		// Action: Update(replace) entire cluster info
		ClusterMap[*cluster.Id] = cluster
		ctx.StatusCode(iris.StatusOK)
	} else {
		// Action: Generate a UUID clusterID and add cluster info.
		clusterId := uuid.New()
		ClusterMap[clusterId.String()] = cluster
		ctx.StatusCode(iris.StatusCreated)
	}
}

func DeleteCluster(ctx iris.Context) {
	clusterId := ctx.Params().Get("clusterId")

	// Validation: Check cluster with this ID exists?
	if !ValidClusterId(clusterId) {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid cluster ID"))
		return
	}

	delete(ClusterMap, clusterId)
	ctx.StatusCode(iris.StatusOK)
}
