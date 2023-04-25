package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"log"
	"main/app/cluster/models"
	"main/app/cluster/validation"
	"main/app/tablemux"
	"main/app/tablemux/handlerimplementation"
	"main/app/utils"
	"main/app/utils/zkerrors"
	"os"
	"px.dev/pxapi"
	"strings"
)

type Details struct {
	Domain string `json:"Domain"`
	Url    string `json:"Url"`
}

type ClusterService interface {
	UpdateCluster(ctx iris.Context, cluster models.ClusterDetails) (int, *zkerrors.ZkError)
	DeleteCluster(ctx iris.Context, clusterId string) (int, *zkerrors.ZkError)
	GetResourceDetails(ctx iris.Context, clusterIdx, action, st, apiKey string) models.PixieResponse
	GetNamespaceList(ctx iris.Context, id, st, apiKey string) models.PixieResponse
	GetServiceDetails(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) models.PixieResponse
	GetPodDetailsTimeSeries(ctx iris.Context, clusterIdx, podName, ns, st, apiKey string) map[string]models.PixieResponse
	GetPodList(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) models.PixieResponse
	GetPxlData(ctx iris.Context, clusterIdx, st, apiKey string) models.PixieResponse
}

type clusterService struct {
	pixie tablemux.PixieRepository
}

func NewClusterService(pixie tablemux.PixieRepository) ClusterService {
	return &clusterService{
		pixie: pixie,
	}
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

func (cs *clusterService) UpdateCluster(ctx iris.Context, cluster models.ClusterDetails) (int, *zkerrors.ZkError) {

	var statusCode int
	if cluster.Id != nil {
		// Validation: Check if provided cluster ID is valid?
		//if !ValidateCluster(*cluster.Id, ctx) {
		//	return
		//}
		// Action: Update(replace) entire cluster info
		models.ClusterMap[*cluster.Id] = cluster
		statusCode = iris.StatusOK
	} else {
		// Action: Generate a UUID clusterID and add cluster info.
		clusterId := uuid.New()
		models.ClusterMap[clusterId.String()] = cluster
		statusCode = iris.StatusCreated
	}
	return statusCode, nil
}

func (cs *clusterService) DeleteCluster(ctx iris.Context, clusterId string) (int, *zkerrors.ZkError) {
	delete(models.ClusterMap, clusterId)
	return iris.StatusOK, nil
}

func (cs *clusterService) GetResourceDetails(ctx iris.Context, clusterIdx, action, st, apiKey string) models.PixieResponse {
	var pxResp models.PixieResponse
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		pxResp.Result = nil
		pxResp.ResultsStats = nil
		pxResp.Error = &e
		return pxResp
	}

	if strings.EqualFold(action, "list") {
		return cs.getServiceDetailsList(ctx, clusterIdx, st, apiKey)
	} else if strings.EqualFold(action, "map") {
		return cs.getServiceDetailsMap(ctx, clusterIdx, st, apiKey)
	}

	e := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST, "unsupported action: "+action)
	pxResp.Result = nil
	pxResp.ResultsStats = nil
	pxResp.Error = &e
	return pxResp
}

func (cs *clusterService) GetNamespaceList(ctx iris.Context, id, st, apiKey string) models.PixieResponse {
	var pxResp models.PixieResponse
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		pxResp.Result = nil
		pxResp.ResultsStats = nil
		pxResp.Error = &e
		return pxResp
	}
	stringListMux := handlerimplementation.New[string]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetNamespaceMethodSignature(st), DataFrameName: "my_first_ns"}
	resultSet, err := cs.pixie.GetPixieData(ctx, stringListMux, tx, id, apiKey, details.Domain)

	if resultSet == nil || err != nil {
		pxResp.Result = nil
		pxResp.ResultsStats = nil
		pxResp.Error = err
	} else {
		pxResp.Result = stringListMux.Table.Values
		pxResp.ResultsStats = resultSet.Stats()
		pxResp.Error = nil
	}

	return pxResp
}

func (cs *clusterService) getServiceDetailsMap(ctx iris.Context, id, st, apiKey string) models.PixieResponse {
	var pxResp models.PixieResponse
	serviceMapMux := handlerimplementation.New[handlerimplementation.ServiceMap]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceMapMethodSignature(st), DataFrameName: "my_first_map"}
	resultSet, err := cs.pixie.GetPixieData(ctx, serviceMapMux, tx, id, apiKey, details.Domain)

	if resultSet == nil || err != nil {
		pxResp.Result = nil
		pxResp.ResultsStats = nil
		pxResp.Error = err
	} else {
		pxResp.Result = serviceMapMux.Table.Values
		pxResp.ResultsStats = resultSet.Stats()
		pxResp.Error = nil
	}

	return pxResp
}

func (cs *clusterService) getServiceDetailsList(ctx iris.Context, id, st, apiKey string) models.PixieResponse {
	var pxResp models.PixieResponse
	serviceListMux := handlerimplementation.New[handlerimplementation.Service]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}
	resultSet, err := cs.pixie.GetPixieData(ctx, serviceListMux, tx, id, apiKey, details.Domain)
	if resultSet == nil || err != nil {
		pxResp.Result = nil
		pxResp.ResultsStats = nil
		pxResp.Error = err
	} else {
		pxResp.Result = serviceListMux.Table.Values
		pxResp.ResultsStats = resultSet.Stats()
		pxResp.Error = nil
	}

	return pxResp
}

func (cs *clusterService) GetServiceDetails(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) models.PixieResponse {
	var pxResp models.PixieResponse
	if !validation.ValidatePxlTime(st) {
		err := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		pxResp.Result = nil
		pxResp.ResultsStats = nil
		pxResp.Error = &err
		return pxResp
	}

	var resultSet *pxapi.ScriptResults
	serviceStatMux := handlerimplementation.New[handlerimplementation.ServiceStat]()
	resultSet, err := cs.pixie.GetPixieData(ctx, serviceStatMux, tablemux.MethodTemplate{MethodSignature: utils.GetServiceDetailsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)

	if resultSet == nil || err != nil {
		pxResp.Result = nil
		pxResp.ResultsStats = nil
		pxResp.Error = err
	} else {
		pxResp.Result = serviceStatMux.Table.Values
		pxResp.ResultsStats = resultSet.Stats()
		pxResp.Error = nil
	}
	return pxResp

	//_ = ctx.JSON(map[string]interface{}{
	//	"results": result,
	//	"stats":   resultSet.Stats(),
	//	"status":  200,
	//})
}

func (cs *clusterService) GetPodDetailsTimeSeries(ctx iris.Context, clusterIdx, podName, ns, st, apiKey string) map[string]models.PixieResponse {
	if !validation.ValidatePxlTime(st) {
		return nil
	}

	// for HTTP Requests and HTTP Errors
	reqAndErrMux := handlerimplementation.New[handlerimplementation.PodDetailsErrAndReq]()
	resultSetErrAndReq, errReqAndErr := cs.pixie.GetPixieData(ctx, reqAndErrMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsForHTTPDataAndErrMethodSignature(st, ns+"/"+podName), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	if errReqAndErr != nil {
		log.Println("pod details err and req, error, ", errReqAndErr.Error)
	}
	resultErrAndReq := reqAndErrMux.Table.Values

	// for HTTP Latency
	latencyMux := handlerimplementation.New[handlerimplementation.PodDetailsLatency]()
	resultSetLatency, errLatency := cs.pixie.GetPixieData(ctx, latencyMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsForHTTPLatencyMethodSignature(st, ns+"/"+podName), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	if errLatency != nil {
		log.Println("pod details latency, error, ", errLatency.Error)
	}
	resultLatency := latencyMux.Table.Values

	// for CPU Usage
	cpuUsageMux := handlerimplementation.New[handlerimplementation.PodDetailsCpuUsage]()
	resultSetCpuUsage, errCpuUsage := cs.pixie.GetPixieData(ctx, cpuUsageMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsForCpuUsageMethodSignature(st, ns+"/"+podName), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	if errCpuUsage != nil {
		log.Println("pod details cpu usage, error, ", errCpuUsage.Error)
	}
	resultCpuUsage := cpuUsageMux.Table.Values

	data := map[string]models.PixieResponse{}
	data["requestAndError"] = models.PixieResponse{
		Result:       resultErrAndReq,
		ResultsStats: resultSetErrAndReq.Stats(),
		Error:        errReqAndErr,
	}
	data["latency"] = models.PixieResponse{
		Result:       resultLatency,
		ResultsStats: resultSetLatency.Stats(),
		Error:        errLatency,
	}
	data["cpuUsage"] = models.PixieResponse{
		Result:       resultCpuUsage,
		ResultsStats: resultSetCpuUsage.Stats(),
		Error:        errCpuUsage,
	}

	return data
}

func (cs *clusterService) GetPodList(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) models.PixieResponse {
	var pxResp models.PixieResponse
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		pxResp.Result = nil
		pxResp.ResultsStats = nil
		pxResp.Error = &e
		return pxResp
	}
	var resultSet *pxapi.ScriptResults

	serviceStatMux := handlerimplementation.New[handlerimplementation.PodDetails]()
	resultSet, err := cs.pixie.GetPixieData(ctx, serviceStatMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsMethodSignature(st, ns, ns+"/"+name), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)

	if resultSet == nil || err != nil {
		pxResp.Result = nil
		pxResp.ResultsStats = nil
		pxResp.Error = err
	} else {
		pxResp.Result = serviceStatMux.Table.Values
		pxResp.ResultsStats = resultSet.Stats()
		pxResp.Error = nil
	}
	return pxResp
	//
	//_ = ctx.JSON(map[string]interface{}{
	//	"results": result,
	//	"stats":   resultSet.Stats(),
	//	"status":  200,
	//})
}

func (cs *clusterService) GetPxlData(ctx iris.Context, clusterIdx, st, apiKey string) models.PixieResponse {
	var pxResp models.PixieResponse
	if !validation.ValidatePxlTime(st) {
		err := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		pxResp.Result = nil
		pxResp.ResultsStats = nil
		pxResp.Error = &err
		return pxResp
	}
	pixieTraceDataMux := handlerimplementation.New[handlerimplementation.PixieTraceData]()

	tx := tablemux.MethodTemplate{MethodSignature: utils.GetPXDataSignature(100, st, "{}"), DataFrameName: "my_first_list"}
	resultSet, err := cs.pixie.GetPixieData(ctx, pixieTraceDataMux, tx, clusterIdx, apiKey, details.Domain)

	if resultSet == nil || err != nil {
		pxResp.Result = nil
		pxResp.ResultsStats = nil
		pxResp.Error = err
	} else {
		pxResp.Result = pixieTraceDataMux.Table.Values
		pxResp.ResultsStats = resultSet.Stats()
		pxResp.Error = nil
	}
	return pxResp
}

func (cs *clusterService) getResp(resultSet *pxapi.ScriptResults, result interface{}) map[string]interface{} {
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

//func (o *MyTestObject) SavePersonDetails(firstname, lastname string, age int) (int, error) {
//	args := o.Called(firstname, lastname, age)
//	return args.Int(0), args.Error(1)
//}
//
//args.Int(0)
//args.Bool(1)
//args.String(2)
//
//return args.Get(0).(*MyObject), args.Get(1).(*AnotherObjectOfMine)
