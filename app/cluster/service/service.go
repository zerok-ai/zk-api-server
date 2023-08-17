package service

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"log"
	"os"
	"px.dev/pxapi"
	"zk-api-server/app/cluster/transformer"
	"zk-api-server/app/cluster/validation"
	"zk-api-server/app/tablemux"
	"zk-api-server/app/tablemux/handlerimplementation"
	"zk-api-server/app/utils"
	"zk-api-server/app/utils/errors"
)

var LogTag = "Service/Service"

type Details struct {
	Domain string `json:"Domain"`
	Url    string `json:"Url"`
}

type ClusterService interface {
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

func (cs *clusterService) GetNamespaceList(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[string], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
		return nil, &e
	}
	mux := handlerimplementation.New[string]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetNamespaceMethodSignature(st), DataFrameName: "my_first_ns"}
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tx, id, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err

}

func (cs *clusterService) GetServiceDetailsMap(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.ServiceMap], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
		return nil, &e
	}

	mux := handlerimplementation.New[handlerimplementation.ServiceMap]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceMapMethodSignature(st), DataFrameName: "my_first_map"}
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tx, id, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err

}

func (cs *clusterService) GetMysqlServiceList(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.ServiceName], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
		return nil, &e
	}

	mux := handlerimplementation.New[handlerimplementation.ServiceName]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetMysqlServiceListMethodSignature(st), DataFrameName: "my_first_list"}
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tx, id, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err
}

func (cs *clusterService) GetPgsqlServiceList(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.ServiceName], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
		return nil, &e
	}

	mux := handlerimplementation.New[handlerimplementation.ServiceName]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetPgsqlServiceListMethodSignature(st), DataFrameName: "my_first_list"}
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tx, id, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err
}

func (cs *clusterService) GetHttpServiceList(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.ServiceName], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
		return nil, &e
	}

	mux := handlerimplementation.New[handlerimplementation.ServiceName]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetHttpServiceListMethodSignature(st), DataFrameName: "my_first_list"}
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tx, id, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err
}

func (cs *clusterService) GetServiceDetailsList(ctx iris.Context, id, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.Service], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
		return nil, &e
	}

	protocolMap := make(map[string]string)

	httpSvcList, _ := cs.GetHttpServiceList(ctx, id, st, apiKey)
	zkLogger.Debug(LogTag, "Http ", httpSvcList.Results)

	for _, svc := range httpSvcList.Results {
		name := svc.Name
		if name != nil {
			protocolMap[*name] = "http"
		}
	}

	mysqlSvcList, _ := cs.GetMysqlServiceList(ctx, id, st, apiKey)
	zkLogger.Debug(LogTag, "Mysql ", mysqlSvcList.Results)

	for _, svc := range mysqlSvcList.Results {
		name := svc.Name
		if name != nil {
			protocolMap[*name] = "mysql"
		}
	}

	pgsqlSvcList, _ := cs.GetPgsqlServiceList(ctx, id, st, apiKey)
	zkLogger.Debug(LogTag, "Pgsql ", pgsqlSvcList.Results)

	for _, svc := range pgsqlSvcList.Results {
		name := svc.Name
		if name != nil {
			protocolMap[*name] = "pgsql"
		}
	}

	mux := handlerimplementation.New[handlerimplementation.Service]()
	tx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tx, id, apiKey, details.Domain)
	response := transformer.PixieResponseToHTTPResponse(resultSet, mux, err)
	results := response.Results

	var newResults []handlerimplementation.Service
	for _, result := range results {
		name := result.ServiceName
		if name != nil {
			protocol, ok := protocolMap[*name]
			if ok {
				result.Protocol = &protocol
			}
		}
		newResults = append(newResults, result)
	}
	response.Results = newResults
	return response, err
}

func (cs *clusterService) GetServiceDetails(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.ServiceStat], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		err := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
		return nil, &err
	}

	var resultSet *pxapi.ScriptResults
	mux := handlerimplementation.New[handlerimplementation.ServiceStat]()
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tablemux.MethodTemplate{MethodSignature: utils.GetServiceDetailsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err

}

func (cs *clusterService) GetPodDetailsTimeSeries(ctx iris.Context, clusterIdx, podName, ns, st, apiKey string) (*transformer.PodDetailsPixieHTTPResponse, *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		err := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
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
		return nil, zkCommon.ToPtr[zkerrors.ZkError](zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil))
	}
	return transformer.PixieResponseToPodDetailsHTTPResponse(reqAndErrHttpResp, httpLatencyHttpResp, cpuUsageHttpResp), nil

}

func (cs *clusterService) GetPodList(ctx iris.Context, clusterIdx, name, ns, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.PodDetails], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		e := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
		return nil, &e
	}
	var resultSet *pxapi.ScriptResults

	mux := handlerimplementation.New[handlerimplementation.PodDetails]()
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsMethodSignature(st, ns, ns+"/"+name), DataFrameName: "my_first_graph"}, clusterIdx, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err
}

func (cs *clusterService) GetPxlData(ctx iris.Context, clusterIdx, st, apiKey string) (*transformer.PixieHTTPResponse[handlerimplementation.PixieTraceData], *zkerrors.ZkError) {
	if !validation.ValidatePxlTime(st) {
		err := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
		return nil, &err
	}
	mux := handlerimplementation.New[handlerimplementation.PixieTraceData]()

	tx := tablemux.MethodTemplate{MethodSignature: utils.GetPXDataSignature(100, st, "{}"), DataFrameName: "my_first_list"}
	resultSet, err := cs.pixie.GetPixieData(ctx, mux, tx, clusterIdx, apiKey, details.Domain)
	return transformer.PixieResponseToHTTPResponse(resultSet, mux, err), err

}
