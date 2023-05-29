package cluster

import (
	"github.com/kataras/iris/v12/core/router"
	"main/app/cluster/handler"
	handler2 "main/app/scenario/handler"
	"main/app/scenario/repository"
	"main/app/scenario/service"
	"main/app/utils"
)

func Initialize(app router.Party) {
	ch := handler.NewClusterHandler()
	{
		clusterAPI := app.Party("/u/cluster")
		clusterAPI.Get("/{clusterIdx}/service/list", utils.ValidateApiKeyMiddleware, ch.GetServiceDetailsList)
		clusterAPI.Get("/{clusterIdx}/service/map", utils.ValidateApiKeyMiddleware, ch.GetServiceDetailsMap)
		clusterAPI.Get("/{clusterIdx}/service/details", utils.ValidateApiKeyMiddleware, ch.GetServiceDetails)
		clusterAPI.Get("/{clusterIdx}/pod/list", utils.ValidateApiKeyMiddleware, ch.GetPodList)
		clusterAPI.Get("/{clusterIdx}/pod/details", utils.ValidateApiKeyMiddleware, ch.GetPodDetails)
		clusterAPI.Get("/traces", utils.ValidateApiKeyMiddleware, ch.GetPxData)
	}

	rr := repository.NewZkPostgresRepo()
	rs := service.NewScenarioService(rr)
	rh := handler2.NewScenarioHandler(rs)
	ruleEngineAPI := app.Party("/o/cluster")
	{
		ruleEngineAPI.Get("/rules", rh.GetAllScenario)
	}
}
