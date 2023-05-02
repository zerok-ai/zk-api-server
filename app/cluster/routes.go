package cluster

import (
	"github.com/kataras/iris/v12/core/router"
	"main/app/cluster/handler"
	handler2 "main/app/ruleengine/handler"
	"main/app/ruleengine/repository"
	"main/app/ruleengine/service"
	"main/app/utils"
)

func Initialize(app router.Party) {
	ch := handler.NewClusterHandler()
	{
		clusterAPI := app.Party("/u/cluster")
		clusterAPI.Post("/", ch.UpsertCluster)
		clusterAPI.Delete("/{clusterId}", utils.ValidateApiKeyMiddleware, ch.DeleteCluster)
		clusterAPI.Get("/{clusterIdx}/service/list", utils.ValidateApiKeyMiddleware, ch.GetResourceDetailsList)
		clusterAPI.Get("/{clusterIdx}/service/map", utils.ValidateApiKeyMiddleware, ch.GetResourceDetailsMap)
		clusterAPI.Get("/{clusterIdx}/service/details", utils.ValidateApiKeyMiddleware, ch.GetServiceDetails)
		clusterAPI.Get("/{clusterIdx}/pod/list", utils.ValidateApiKeyMiddleware, ch.GetPodList)
		clusterAPI.Get("/{clusterIdx}/pod/details", utils.ValidateApiKeyMiddleware, ch.GetPodDetails)
		clusterAPI.Get("/traces", utils.ValidateApiKeyMiddleware, ch.GetPxData)
	}

	rr := repository.NewRulesFromFileRepo()
	rs := service.NewRuleService(rr)
	rh := handler2.NewRuleHandler(rs)
	ruleEngineAPI := app.Party("/o/cluster")
	{
		ruleEngineAPI.Get("/rules", rh.GetAllRules)
	}
}
