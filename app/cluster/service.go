package cluster

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"main/app/cluster/models"
	"main/app/tablemux"
	"main/app/tablemux/handlerimplementation"
	"main/app/utils"
	"px.dev/pxapi"
	"strings"
)

func init() {
	tablemux.GetAuthTokenWith2ReTry(3)
	tablemux.PopulateApiKey()
}

func listCluster(ctx iris.Context) {

	tablemux.GetAuthTokenWith2ReTry(3)
	r := tablemux.GetMetaDataWithRetry(3)

	if r.StatusCode == 200 {
		finalResp := handlerimplementation.ClusterDetailsMetaDataResponse{}
		responseData, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseData, &finalResp)
		tablemux.UpdateApiKey(finalResp.Data.ApiKey.Key)
		clusters := make([]models.ClusterDetails, 0)

		for _, r := range finalResp.Data.Clusters {
			clusters = append(clusters, models.FromResponseToDomainClusterDetails(r))
		}

		err := ctx.JSON(clusters)
		if err != nil {
			return
		}
	}
	ctx.StatusCode(iris.StatusInternalServerError)
	ctx.SetErr(utils.ErrClusterFetchFailed)
	return
}

func updateCluster(ctx iris.Context, cluster models.ClusterDetails) {

	// Create a cluster if cluster.ID is missing. else, update if valid.
	if cluster.Id != nil {
		// Validation: Check if provided cluster ID is valid?
		//if !ValidateCluster(*cluster.Id, ctx) {
		//	return
		//}
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
	//if !ValidateCluster(clusterId, ctx) {
	//	return
	//}

	delete(models.ClusterMap, clusterId)
	ctx.StatusCode(iris.StatusOK)
}

func getResourceDetails(ctx iris.Context, clusterIdx, action, st string) {
	if !ValidatePxlTime(ctx, st) {
		return
	}

	var resultSet *pxapi.ScriptResults
	var result interface{}
	if strings.EqualFold(action, "list") {
		resultSet, result = getServiceDetailsList(ctx, clusterIdx, st)
	} else if strings.EqualFold(action, "map") {
		resultSet, result = getServiceDetailsMap(ctx, clusterIdx, st)
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

func getNamespaceList(ctx iris.Context, id, st string) (*pxapi.ScriptResults, []string) {
	var v = make([]string, 0)
	stringListMux := handlerimplementation.StringListMux{Table: handlerimplementation.TablePrinterStringList{Values: v}}
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetNamespaceMethodSignature(st), DataFrameName: "my_first_ns"}
	resultSet := tablemux.GetResource(ctx, id, &stringListMux, tx, 3)
	return resultSet, stringListMux.Table.Values
}

func getServiceDetailsMap(ctx iris.Context, id, st string) (*pxapi.ScriptResults, []handlerimplementation.ServiceMap) {
	var s = make([]handlerimplementation.ServiceMap, 0)
	serviceListMux := handlerimplementation.ServiceMapMux{Table: handlerimplementation.TablePrinterServiceMap{Values: s}}
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceMapMethodSignature(st), DataFrameName: "my_first_map"}
	resultSet := tablemux.GetResource(ctx, id, &serviceListMux, tx, 3)
	return resultSet, serviceListMux.Table.Values
}

func getServiceDetailsList(ctx iris.Context, id, st string) (*pxapi.ScriptResults, []handlerimplementation.Service) {
	var s = make([]handlerimplementation.Service, 0)
	serviceListMux := handlerimplementation.ServiceListMux{Table: handlerimplementation.TablePrinterServiceList{Values: s}}
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}
	resultSet := tablemux.GetResource(ctx, id, &serviceListMux, tx, 3)
	return resultSet, serviceListMux.Table.Values
}

func getServiceDetails(ctx iris.Context, clusterIdx, name, ns, st string) {
	if !ValidatePxlTime(ctx, st) {
		return
	}
	var resultSet *pxapi.ScriptResults
	var result interface{}

	var s = make([]handlerimplementation.ServiceStat, 0)
	serviceStatMux := handlerimplementation.ServiceStatMux{Table: handlerimplementation.TablePrinterServiceStat{Values: s}}
	resultSet = tablemux.GetResource(ctx, clusterIdx, &serviceStatMux, tablemux.MethodTemplate{MethodSignature: utils.GetServiceDetailsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"}, 3)
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
