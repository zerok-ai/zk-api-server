package app

import (
	"github.com/kataras/iris/v12/core/router"
	"zk-api-server/app/attribute/handler"
	clusterHandler "zk-api-server/app/cluster/handler"
	integrationsHandler "zk-api-server/app/integrations/handler"
	scenarioHandler "zk-api-server/app/scenario/handler"
	"zk-api-server/app/utils"
)

func Initialize(app router.Party, rh scenarioHandler.ScenarioHandler, ch clusterHandler.ClusterHandler, ih integrationsHandler.IntegrationsHandler, ah handler.AttributeHandler) {
	{
		clusterAPI := app.Party("/u/cluster")
		clusterAPI.Get("/{"+utils.ClusterIdxPathParam+"}/service/list", utils.ValidateApiKeyMiddleware, ch.GetServiceDetailsList)
		clusterAPI.Get("/{"+utils.ClusterIdxPathParam+"}/service/map", utils.ValidateApiKeyMiddleware, ch.GetServiceDetailsMap)
		clusterAPI.Get("/{"+utils.ClusterIdxPathParam+"}/service/details", utils.ValidateApiKeyMiddleware, ch.GetServiceDetails)
		clusterAPI.Get("/{"+utils.ClusterIdxPathParam+"}/pod/list", utils.ValidateApiKeyMiddleware, ch.GetPodList)
		clusterAPI.Get("/{"+utils.ClusterIdxPathParam+"}/pod/details", utils.ValidateApiKeyMiddleware, ch.GetPodDetails)
		clusterAPI.Get("/{"+utils.ClusterIdxPathParam+"}/traces", utils.ValidateApiKeyMiddleware, ch.GetPxData)
		clusterAPI.Get("/{"+utils.ClusterIdxPathParam+"}/scenario", rh.GetAllScenarioDashboard)
		clusterAPI.Post("/{"+utils.ClusterIdxPathParam+"}/scenario", rh.CreateScenario)
		clusterAPI.Put("/{"+utils.ClusterIdxPathParam+"}/scenario/{scenarioIdx}/status", rh.UpdateScenarioState)
		clusterAPI.Delete("/{"+utils.ClusterIdxPathParam+"}/scenario/{scenarioIdx}", rh.DeleteScenario)

		clusterAPI.Post("/{"+utils.ClusterIdxPathParam+"}/integration", ih.UpsertIntegration)
		clusterAPI.Get("/{"+utils.ClusterIdxPathParam+"}/integration", ih.GetAllIntegrationsDashboard)

		clusterAPI.Get("/attribute", ah.GetAttributes)
		clusterAPI.Put("/attribute", ah.UploadAttributesCSV)

	}

	ruleEngineAPI := app.Party("/o/cluster")
	{
		ruleEngineAPI.Get("/{"+utils.ClusterIdxPathParam+"}/scenario", rh.GetAllScenarioOperator)

		ruleEngineAPI.Get("/{"+utils.ClusterIdxPathParam+"}/integration", ih.GetAllIntegrationsOperator)

		ruleEngineAPI.Get("/attribute", ah.GetAttributesForBackend)

	}
}
