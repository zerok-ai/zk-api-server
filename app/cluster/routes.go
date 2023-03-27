package cluster

import (
	"github.com/kataras/iris/v12/core/router"
)

func Initialize(app router.Party) {
	clusterAPI := app.Party("/u/cluster")
	{
		clusterAPI.Post("/", UpsertCluster)
		clusterAPI.Delete("/{clusterId}", DeleteCluster)
		clusterAPI.Get("/{clusterIdx}/service/list", GetResourceDetailsList)
		clusterAPI.Get("/{clusterIdx}/service/map", GetResourceDetailsMap)
		clusterAPI.Get("/{clusterIdx}/service/details", GetServiceDetails)
		clusterAPI.Get("/{clusterIdx}/pod/list", GetPodDetailsList)
		clusterAPI.Get("/traces", GetPxData)
	}
}
