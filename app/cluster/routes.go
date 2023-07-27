package cluster

import (
	"github.com/kataras/iris/v12/core/router"
	clusterHandler "main/app/cluster/handler"
	scenarioHandler "main/app/scenario/handler"
	"main/app/utils"
)

func Initialize(app router.Party, rh scenarioHandler.ScenarioHandler, ch clusterHandler.ClusterHandler) {
	{
		clusterAPI := app.Party("/u/cluster")
		clusterAPI.Get("/{clusterIdx}/service/list", utils.ValidateApiKeyMiddleware, ch.GetServiceDetailsList)
		clusterAPI.Get("/{clusterIdx}/service/map", utils.ValidateApiKeyMiddleware, ch.GetServiceDetailsMap)
		clusterAPI.Get("/{clusterIdx}/service/details", utils.ValidateApiKeyMiddleware, ch.GetServiceDetails)
		clusterAPI.Get("/{clusterIdx}/pod/list", utils.ValidateApiKeyMiddleware, ch.GetPodList)
		clusterAPI.Get("/{clusterIdx}/pod/details", utils.ValidateApiKeyMiddleware, ch.GetPodDetails)
		clusterAPI.Get("/traces", utils.ValidateApiKeyMiddleware, ch.GetPxData)
		clusterAPI.Get("/scenario", rh.GetAllScenario)
	}

	ruleEngineAPI := app.Party("/o/cluster")
	{
		ruleEngineAPI.Get("/scenario", rh.GetAllScenario)
	}
}
