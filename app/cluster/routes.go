package cluster

import (
	"github.com/kataras/iris/v12/core/router"
	"main/app/cluster/handler"
)

func Initialize(app router.Party) {

	ch := handler.NewClusterHandler()
	clusterAPI := app.Party("/u/cluster")
	{
		clusterAPI.Post("/", ch.UpsertCluster)
		clusterAPI.Delete("/{clusterId}", ch.DeleteCluster)
		clusterAPI.Get("/{clusterIdx}/service/list", ch.GetResourceDetailsList)
		clusterAPI.Get("/{clusterIdx}/service/map", ch.GetResourceDetailsMap)
		clusterAPI.Get("/{clusterIdx}/service/details", ch.GetServiceDetails)
		clusterAPI.Get("/{clusterIdx}/pod/list", ch.GetPodList)
		clusterAPI.Get("/{clusterIdx}/pod/details", ch.GetPodDetails)
		clusterAPI.Get("/traces", ch.GetPxData)
	}
}
