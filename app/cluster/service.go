package cluster

import (
	"encoding/json"
	"log"
	"main/app/cluster/models"
	"main/app/tablemux"
	"main/app/tablemux/handlerimplementation"
	"main/app/utils"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"px.dev/pxapi"
)

type Details struct {
	Domain string `json:"Domain"`
	Url    string `json:"Url"`
}

var details Details

func init() {
	configFilePath := "cluster.conf"

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

	println(resultSet)

	_ = ctx.JSON(map[string]interface{}{
		"results": result,
		"stats":   resultSet.Stats(),
		"status":  200,
	})
}

func getNamespaceList(ctx iris.Context, id, st, apiKey string) (*pxapi.ScriptResults, []string) {
	stringListMux := handlerimplementation.New[string]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetNamespaceMethodSignature(st), DataFrameName: "my_first_ns"}
	resultSet := tablemux.GetResource(ctx, stringListMux, tx, id, apiKey, details.Domain)
	return resultSet, stringListMux.Table.Values
}

func getServiceDetailsMap(ctx iris.Context, id, st, apiKey string) (*pxapi.ScriptResults, []handlerimplementation.ServiceMap) {
	serviceMapMux := handlerimplementation.New[handlerimplementation.ServiceMap]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceMapMethodSignature(st), DataFrameName: "my_first_map"}
	resultSet := tablemux.GetResource(ctx, serviceMapMux, tx, id, apiKey, details.Domain)
	return resultSet, serviceMapMux.Table.Values
}

func getServiceDetailsList(ctx iris.Context, id, st, apiKey string) (*pxapi.ScriptResults, []handlerimplementation.Service) {
	serviceListMux := handlerimplementation.New[handlerimplementation.Service]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}
	resultSet := tablemux.GetResource(ctx, serviceListMux, tx, id, apiKey, details.Domain)
	return resultSet, serviceListMux.Table.Values
}

func getServiceDetails(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) {
	if !ValidatePxlTime(ctx, st) {
		return
	}
	var resultSet *pxapi.ScriptResults
	var result interface{}

	serviceStatMux := handlerimplementation.New[handlerimplementation.ServiceState]()
	resultSet = tablemux.GetResource(ctx, serviceStatMux, tablemux.MethodTemplate{MethodSignature: utils.GetServiceDetailsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
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

func getPodDetailsTimeSeries(ctx iris.Context, clusterIdx, podName, ns, st, apiKey string) {
	if !ValidatePxlTime(ctx, st) {
		return
	}

	// for HTTP Requests and HTTP Errors
	reqAndErrMux := handlerimplementation.New[handlerimplementation.PodDetailsErrAndReq]()
	resultSetErrAndReq := tablemux.GetResource(ctx, reqAndErrMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsForHTTPDataAndErrMethodSignature(st, ns+"/"+podName), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	resultErrAndReq := reqAndErrMux.Table.Values

	// for HTTP Latency
	latencyMux := handlerimplementation.New[handlerimplementation.PodDetailsLatency]()
	resultSetLatency := tablemux.GetResource(ctx, latencyMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsForHTTPLatencyMethodSignature(st, ns+"/"+podName), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	resultLatency := latencyMux.Table.Values

	// for CPU Usage
	cpuUsageMux := handlerimplementation.New[handlerimplementation.PodDetailsCpuUsage]()
	resultSetCpuUsage := tablemux.GetResource(ctx, cpuUsageMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsForCpuUsageMethodSignature(st, ns+"/"+podName), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	resultCpuUsage := cpuUsageMux.Table.Values

	reqAndErrResp := getResp(resultSetErrAndReq, resultErrAndReq)
	latencyResp := getResp(resultSetLatency, resultLatency)
	cpuUsageResp := getResp(resultSetCpuUsage, resultCpuUsage)

	_ = ctx.JSON(map[string]map[string]interface{}{
		"errAndReq": reqAndErrResp,
		"latency":   latencyResp,
		"cpuUsage":  cpuUsageResp,
	})
}

func getPodDetails(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) {
	if !ValidatePxlTime(ctx, st) {
		return
	}
	var resultSet *pxapi.ScriptResults
	var result interface{}

	serviceStatMux := handlerimplementation.New[handlerimplementation.PodDetails]()
	resultSet = tablemux.GetResource(ctx, serviceStatMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsMethodSignature(st, ns, ns+"/"+name), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
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

	pixieTraceDataMux := handlerimplementation.New[handlerimplementation.PixieTraceData]()

	tx := tablemux.MethodTemplate{MethodSignature: utils.GetPXDataSignature(100, st, "{}"), DataFrameName: "my_first_list"}
	resultSet := tablemux.GetResource(ctx, pixieTraceDataMux, tx, clusterIdx, apiKey, details.Domain)
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

func getResp(resultSet *pxapi.ScriptResults, result interface{}) map[string]interface{} {
	var x map[string]interface{}
	if result == nil {
		x = map[string]interface{}{
			"results": nil,
			"stats":   nil,
			"status":  500,
		}
	} else {
		x = map[string]interface{}{
			"results": result,
			"stats":   resultSet.Stats(),
			"status":  200,
		}
	}
	return x
}
