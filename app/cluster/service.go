package cluster

import (
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"main/app/cluster/models"
	"main/app/cluster/models/tablemux"
	"main/app/utils"
	"px.dev/pxapi"
	"strings"
)

func listCluster(ctx iris.Context) {
	var clusters []models.Cluster
	for k, v := range models.ClusterMap {
		var id = k
		v.Id = &id
		clusters = append(clusters, v)
	}
	if clusters == nil {
		var emptyArr []string
		err := ctx.JSON(emptyArr)
		if err != nil {
			return
		}
		return
	}
	err := ctx.JSON(clusters)
	if err != nil {
		return
	}
}

func updateCluster(ctx iris.Context, cluster models.Cluster) {

	// Create a cluster if cluster.ID is missing. else, update if valid.
	if cluster.Id != nil {
		// Validation: Check if provided cluster ID is valid?
		if !ValidateCluster(*cluster.Id, ctx) {
			return
		}
		// Action: Update(replace) entire cluster info
		models.ClusterMap[*cluster.Id] = cluster
		ctx.StatusCode(iris.StatusOK)
	} else {
		// Action: Generate a UUID clusterID and add cluster info.
		clusterId := uuid.New()
		models.ClusterMap[clusterId.String()] = cluster
		ctx.StatusCode(iris.StatusCreated)
	}
}

func deleteCluster(ctx iris.Context, clusterId string) {
	if !ValidateCluster(clusterId, ctx) {
		return
	}

	delete(models.ClusterMap, clusterId)
	ctx.StatusCode(iris.StatusOK)
}

func getResourceDetails(ctx iris.Context, clusterIdx, namespace, resource, action, st string) {

	if !ValidateResource(resource, ctx) {
		return
	}
	if !ValidateAction(action, ctx) {
		return
	}
	if !ValidateCluster(clusterIdx, ctx) {
		return
	}

	clusterDetails := models.ClusterMap[clusterIdx]
	apiKey := clusterDetails.ApiKey
	cloudAddress := clusterDetails.Domain + ":443"
	clusterId := clusterDetails.ClusterId

	if namespace != "all" {
		var v = make([]string, 0)
		stringListMux := tablemux.StringListMux{Table: tablemux.TablePrinterStringList{Values: v}}

		tx := models.MethodTemplate{MethodSignature: utils.GetNamespaceMethodSignature(st), DataFrameName: "my_first_ns"}
		_ = getResource(ctx, apiKey, clusterId, cloudAddress, &stringListMux, tx)
		namespaceList := stringListMux.Table.Values
		if !utils.Contains(namespaceList, namespace) {
			_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
				Title("Invalid namespace"))
			return
		}
	}

	var resultSet *pxapi.ScriptResults
	var result interface{}

	//tx := models.Template{StartTime: "-20m", Resource: "pod"}
	if strings.EqualFold(action, "list") {
		var s = make([]models.Service, 0)
		serviceListMux := tablemux.ServiceListMux{Table: tablemux.TablePrinterServiceList{Values: s}}
		tx := models.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}
		resultSet = getResource(ctx, apiKey, clusterId, cloudAddress, &serviceListMux, tx)
		result = serviceListMux.Table.Values
	} else if strings.EqualFold(action, "map") {
		var s = make([]models.ServiceMap, 0)
		serviceListMux := tablemux.ServiceMapMux{Table: tablemux.TablePrinterServiceMap{Values: s}}
		tx := models.MethodTemplate{MethodSignature: utils.GetServiceMapMethodSignature(st), DataFrameName: "my_first_map"}
		resultSet = getResource(ctx, apiKey, clusterId, cloudAddress, &serviceListMux, tx)
		result = serviceListMux.Table.Values
	}

	if result == nil {
		return
	}

	_ = ctx.JSON(map[string]interface{}{
		"results": result,
		"stats":   resultSet.Stats(),
		"status":  200,
	})
}

func getResource(ctx iris.Context, apiKey, clusterId, cloudAddress string, t tablemux.TableRecordHandler, tx models.MethodTemplate) *pxapi.ScriptResults {
	vz, pxl, ctxNew, err := tablemux.A(apiKey, clusterId, cloudAddress, tx)
	if err != nil {
		_ = ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title(err.Error()))
		return nil
	}
	resultSet, err := t.B(ctxNew, vz, pxl)
	if err != nil {
		_ = ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title(err.Error()))
		return nil
	}
	resultSet, err = tablemux.C(resultSet)
	if err != nil {
		_ = ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title(err.Error()))
		return nil
	}
	return resultSet
}

func getServiceStatsGraph(ctx iris.Context, clusterIdx, name, ns, st string) {
	if !ValidateCluster(clusterIdx, ctx) {
		return
	}

	//TODO: write validation for st

	clusterDetails := models.ClusterMap[clusterIdx]
	apiKey := clusterDetails.ApiKey
	cloudAddress := clusterDetails.Domain + ":443"
	clusterId := clusterDetails.ClusterId

	var resultSet *pxapi.ScriptResults
	var result interface{}

	var s = make([]models.ServiceStat, 0)
	serviceStatMux := tablemux.ServiceStatMux{Table: tablemux.TablePrinterServiceStat{Values: s}}
	resultSet = getResource(ctx, apiKey, clusterId, cloudAddress, &serviceStatMux, models.MethodTemplate{MethodSignature: utils.GetServiceStatsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"})
	result = serviceStatMux.Table.Values

	if result == nil {
		return
	}

	_ = ctx.JSON(map[string]interface{}{
		"results": result,
		"stats":   resultSet.Stats(),
		"status":  200,
	})
}
