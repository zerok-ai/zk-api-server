package cluster

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"log"
	"main/app/cluster/models"
	"main/app/tablemux"
	"main/app/tablemux/handlerimplementation"
	"main/app/utils"
	"os"
	"px.dev/pxapi"
	"strings"
)

type Details struct {
	Domain string `json:"Domain"`
	Url    string `json:"Url"`
}

var details Details

func init() {
	configFilePath := "/opt/cluster.conf"

	jsonFile, err := os.Open(configFilePath)

	if err != nil {
		log.Println(err)
		os.Exit(2)
		return
	} else {
		defer jsonFile.Close()

		err = json.NewDecoder(jsonFile).Decode(&details)
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
	}
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

func getResourceDetails(ctx iris.Context, clusterIdx, action, st, apiKey string) {
	if !ValidatePxlTime(ctx, st) {
		return
	}

	var resultSet *pxapi.ScriptResults
	var result interface{}
	if strings.EqualFold(action, "list") {
		resultSet, result = getServiceDetailsList(ctx, clusterIdx, st, apiKey)
	} else if strings.EqualFold(action, "map") {
		resultSet, result = getServiceDetailsMap(ctx, clusterIdx, st, apiKey)
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

func getNamespaceList(ctx iris.Context, id, st, apiKey string) (*pxapi.ScriptResults, []string) {
	var v = make([]string, 0)
	stringListMux := handlerimplementation.StringListMux{Table: handlerimplementation.TablePrinterStringList{Values: v}}
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetNamespaceMethodSignature(st), DataFrameName: "my_first_ns"}
	resultSet := tablemux.GetResource(ctx, &stringListMux, tx, id, apiKey, details.Domain)
	return resultSet, stringListMux.Table.Values
}

func getServiceDetailsMap(ctx iris.Context, id, st, apiKey string) (*pxapi.ScriptResults, []handlerimplementation.ServiceMap) {
	var s = make([]handlerimplementation.ServiceMap, 0)
	serviceMapMux := handlerimplementation.ServiceMapMux{Table: handlerimplementation.TablePrinterServiceMap{Values: s}}
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceMapMethodSignature(st), DataFrameName: "my_first_map"}
	resultSet := tablemux.GetResource(ctx, &serviceMapMux, tx, id, apiKey, details.Domain)
	return resultSet, serviceMapMux.Table.Values
}

func getServiceDetailsList(ctx iris.Context, id, st, apiKey string) (*pxapi.ScriptResults, []handlerimplementation.Service) {
	var s = make([]handlerimplementation.Service, 0)
	serviceListMux := handlerimplementation.ServiceListMux{Table: handlerimplementation.TablePrinterServiceList{Values: s}}
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}
	resultSet := tablemux.GetResource(ctx, &serviceListMux, tx, id, apiKey, details.Domain)
	return resultSet, serviceListMux.Table.Values
}

func getServiceDetails(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) {
	if !ValidatePxlTime(ctx, st) {
		return
	}
	var resultSet *pxapi.ScriptResults
	var result interface{}

	var s = make([]handlerimplementation.ServiceStat, 0)
	serviceStatMux := handlerimplementation.ServiceStatMux{Table: handlerimplementation.TablePrinterServiceStat{Values: s}}
	resultSet = tablemux.GetResource(ctx, &serviceStatMux, tablemux.MethodTemplate{MethodSignature: utils.GetServiceDetailsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
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

func getPodDetails(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) {
	if !ValidatePxlTime(ctx, st) {
		return
	}
	var resultSet *pxapi.ScriptResults
	var result interface{}

	var s = make([]handlerimplementation.PodDetails, 0)
	serviceStatMux := handlerimplementation.PodDetailsListMux{Table: handlerimplementation.TablePrinterPodDetailsList{Values: s}}
	resultSet = tablemux.GetResource(ctx, &serviceStatMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsMethodSignature(st, ns, ns+"/"+name), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
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

func getPxlData(ctx iris.Context, clusterIdx, st, apiKey string) {

	var s = make([]handlerimplementation.PixieTraceData, 0)
	pixieTraceDataMux := handlerimplementation.PixieTraceDataListMux{Table: handlerimplementation.TablePrinterPixieTraceDataList{Values: s}}

	tx := tablemux.MethodTemplate{MethodSignature: utils.GetPXDataSignature(100, st, "{}"), DataFrameName: "my_first_list"}
	resultSet := tablemux.GetResource(ctx, &pixieTraceDataMux, tx, clusterIdx, apiKey, details.Domain)
	result := pixieTraceDataMux.Table.Values

	if result == nil {
		return
	}

	ctx.StatusCode(iris.StatusOK)
	_ = ctx.JSON(map[string]interface{}{
		"results": result,
		"stats":   resultSet.Stats(),
		"status":  200,
	})
}
