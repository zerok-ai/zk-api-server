package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"log"
	"main/app/cluster/models"
	"main/app/cluster/transformer"
	"main/app/cluster/validation"
	"main/app/tablemux"
	"main/app/tablemux/handlerimplementation"
	"main/app/utils"
	"main/app/utils/zkerrors"
	"os"
	"px.dev/pxapi"
)

type Details struct {
	Domain string `json:"Domain"`
	Url    string `json:"Url"`
}

type ClusterService interface {
	UpdateCluster(ctx iris.Context, cluster models.ClusterDetails) (int, *zkerrors.ZkError)
	DeleteCluster(ctx iris.Context, clusterId string) (int, *zkerrors.ZkError)
	GetServiceDetailsMap(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.ServiceMap], *zkerrors.ZkError)
	GetServiceDetailsList(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.Service], *zkerrors.ZkError)
	GetNamespaceList(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[string], *zkerrors.ZkError)
	GetServiceDetails(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.ServiceStat], *zkerrors.ZkError)
	GetPodDetailsTimeSeries(ctx iris.Context, clusterIdx, podName, ns, st, apiKey string) (*transformer.PodDetailsPixieHTTPResponse, *zkerrors.ZkError)
	GetPxlData(ctx iris.Context, clusterIdx, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.PixieTraceData], *zkerrors.ZkError)
	GetPodList(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.PodDetails], *zkerrors.ZkError)
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
	configFilePath := "/Users/vaibhavpaharia/Go/src/zk-api-server/cluster.conf"

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

func (cs *clusterService) GetNamespaceList(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[string], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		return nil, &e
	}
	mux := handlerimplementation.New[string]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetNamespaceMethodSignature(st), DataFrameName: "my_first_ns"}
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tx, id, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err

}

func (cs *clusterService) GetServiceDetailsMap(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.ServiceMap], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		return nil, &e
	}

	mux := handlerimplementation.New[handlerimplementation.ServiceMap]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceMapMethodSignature(st), DataFrameName: "my_first_map"}
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tx, id, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err

}

func (cs *clusterService) GetServiceDetailsList(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.Service], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		return nil, &e
	}

	mux := handlerimplementation.New[handlerimplementation.Service]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tx, id, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err

}

func (cs *clusterService) GetServiceDetails(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.ServiceStat], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		err := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		return nil, &err
	}

	var resultSet *pxapi.ScriptResults
	mux := handlerimplementation.New[handlerimplementation.ServiceStat]()
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tablemux.MethodTemplate{MethodSignature: utils.GetServiceDetailsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err

}

func (cs *clusterService) GetPodDetailsTimeSeries(ctx iris.Context, clusterIdx, podName, ns, st, apiKey string) (*transformer.PodDetailsPixieHTTPResponse, *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		err := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		return nil, &err
	}

	// for HTTP Requests and HTTP Errors
	reqAndErrMux := handlerimplementation.New[handlerimplementation.PodDetailsErrAndReq]()
	resultSetErrAndReq, errReqAndErr := cs.pixie.GetPixieData(ctx, reqAndErrMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsForHTTPDataAndErrMethodSignature(st, ns+"/"+podName), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	reqAndErrHttpResp := transformer.PixieResponseToHTTPResponse(resultSetErrAndReq, reqAndErrMux, errReqAndErr)

	// for HTTP Latency
	latencyMux := handlerimplementation.New[handlerimplementation.PodDetailsLatency]()
	resultSetLatency, errLatency := cs.pixie.GetPixieData(ctx, latencyMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsForHTTPLatencyMethodSignature(st, ns+"/"+podName), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	httpLatencyHttpResp := transformer.PixieResponseToHTTPResponse(resultSetLatency, latencyMux, errLatency)

	// for CPU Usage
	cpuUsageMux := handlerimplementation.New[handlerimplementation.PodDetailsCpuUsage]()
	resultSetCpuUsage, errCpuUsage := cs.pixie.GetPixieData(ctx, cpuUsageMux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsForCpuUsageMethodSignature(st, ns+"/"+podName), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	cpuUsageHttpResp := transformer.PixieResponseToHTTPResponse(resultSetCpuUsage, cpuUsageMux, errCpuUsage)

	if errReqAndErr != nil && errLatency != nil && errCpuUsage != nil {
		return nil, utils.ToPtr[zkerrors.ZkError](zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil))
	}
	return transformer.PixieResponseToPodDetailsHTTPResponse(reqAndErrHttpResp, httpLatencyHttpResp, cpuUsageHttpResp), nil

}

func (cs *clusterService) GetPodList(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.PodDetails], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		return nil, &e
	}
	var resultSet *pxapi.ScriptResults

	mux := handlerimplementation.New[handlerimplementation.PodDetails]()
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsMethodSignature(st, ns, ns+"/"+name), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err
}

func (cs *clusterService) GetPxlData(ctx iris.Context, clusterIdx, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.PixieTraceData], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		err := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
		return nil, &err
	}
	mux := handlerimplementation.New[handlerimplementation.PixieTraceData]()

	tx := tablemux.MethodTemplate{MethodSignature: utils.GetPXDataSignature(100, st, "{}"), DataFrameName: "my_first_list"}
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tx, clusterIdx, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err

}
