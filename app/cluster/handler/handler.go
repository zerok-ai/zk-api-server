package handler

import (
	"github.com/kataras/iris/v12"
	"main/app/cluster/models"
	"main/app/cluster/service"
	"main/app/cluster/validation"
	"main/app/tablemux"
	"main/app/utils"
	"main/app/utils/zkerrors"
	"px.dev/pxapi"
)

type ClusterHandler interface {
	UpsertCluster(ctx iris.Context)
	DeleteCluster(ctx iris.Context)
	GetResourceDetailsList(ctx iris.Context)
	GetResourceDetailsMap(ctx iris.Context)
	GetServiceDetails(ctx iris.Context)
	GetPodList(ctx iris.Context)
	GetPodDetails(ctx iris.Context)
	GetPxData(ctx iris.Context)
}

type clusterHandler struct {
}

func NewClusterHandler() ClusterHandler {
	return &clusterHandler{}
}

var s = service.NewClusterService(tablemux.NewPixieRepository())

func (h *clusterHandler) UpsertCluster(ctx iris.Context) {
	var cluster models.ClusterDetails
	err := ctx.ReadJSON(&cluster)
	// Validation: Is cluster model valid, parse cluster obj
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrClusterParsingFailed)
		return
	}

	statusCode, zkError := s.UpdateCluster(ctx, cluster)
	var zkHttpResponse utils.ZkHttpResponse
	if zkError != nil {
		zkHttpResponse = utils.ZkHttpResponseBuilder{}.WithZkErrorType(zkError.Error).
			Build()
	} else {
		zkHttpResponse = utils.CreateSuccessResponseWithStatusCode(nil, statusCode)
	}
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (h *clusterHandler) DeleteCluster(ctx iris.Context) {
	clusterId := ctx.Params().Get("clusterId")

	if utils.IsEmpty(clusterId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.SetErr(utils.ErrClusterIdEmpty)
		return
	}

	statusCode, zkError := s.DeleteCluster(ctx, clusterId)
	var zkHttpResponse utils.ZkHttpResponse
	if zkError != nil {
		zkHttpResponse = utils.ZkHttpResponseBuilder{}.WithZkErrorType(zkError.Error).
			Build()
	} else {
		zkHttpResponse = utils.CreateSuccessResponseWithStatusCode(nil, statusCode)
	}
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (h *clusterHandler) GetResourceDetailsList(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParam("st")

	if err := validation.ValidateGetResourceDetailsApi(st, apiKey); err != nil {
		zkHttpResponse := utils.CreateErrorResponseWithStatusCode(err.Error)
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	pxResp := s.GetResourceDetails(ctx, clusterIdx, "list", st, apiKey)
	utils.GenerateResponseAndReturn(ctx, pxResp)

}

func (h *clusterHandler) GetResourceDetailsMap(ctx iris.Context) {
	apiKey := ctx.GetHeader("ZK_API_KEY")
	clusterIdx := ctx.Params().Get("clusterIdx")
	st := ctx.URLParam("st")

	if err := validation.ValidateGetResourceDetailsApi(st, apiKey); err != nil {
		zkHttpResponse := utils.CreateErrorResponseWithStatusCode(err.Error)
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	pxResp := s.GetResourceDetails(ctx, clusterIdx, "map", st, apiKey)
	utils.GenerateResponseAndReturn(ctx, pxResp)

}

func (h *clusterHandler) GetServiceDetails(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	apiKey := ctx.GetHeader("ZK_API_KEY")
	serviceName := ctx.URLParam("name")
	ns := ctx.URLParam("ns")
	st := ctx.URLParam("st")

	if err := validation.ValidateGraphDetailsApi(serviceName, ns, st, apiKey); err != nil {
		zkHttpResponse := utils.CreateErrorResponseWithStatusCode(err.Error)
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	pxResp := s.GetServiceDetails(ctx, clusterIdx, serviceName, ns, st, apiKey)
	utils.GenerateResponseAndReturn(ctx, pxResp)

}

func (h *clusterHandler) GetPodList(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParam("st")
	serviceName := ctx.URLParam("service_name")
	ns := ctx.URLParam("ns")

	if err := validation.ValidateGraphDetailsApi(serviceName, ns, st, apiKey); err != nil {
		zkHttpResponse := utils.CreateErrorResponseWithStatusCode(err.Error)
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	pxResp := s.GetPodList(ctx, clusterIdx, serviceName, ns, st, apiKey)
	utils.GenerateResponseAndReturn(ctx, pxResp)

}

func (h *clusterHandler) GetPodDetails(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParam("st")
	podName := ctx.URLParam("pod_name")
	ns := ctx.URLParam("ns")

	if err := validation.ValidatePodDetailsApi(ctx, podName, ns, st, apiKey); err != nil {
		zkHttpResponse := utils.CreateErrorResponseWithStatusCode(err.Error)
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}
	pxRespMap := s.GetPodDetailsTimeSeries(ctx, clusterIdx, podName, ns, st, apiKey)

	reqAndErrResp := getResp(pxRespMap["requestAndError"].ResultsStats, pxRespMap["requestAndError"].Result, pxRespMap["requestAndError"].Error)
	latencyResp := getResp(pxRespMap["latency"].ResultsStats, pxRespMap["latency"].Result, pxRespMap["latency"].Error)
	cpuUsageResp := getResp(pxRespMap["cpuUsage"].ResultsStats, pxRespMap["cpuUsage"].Result, pxRespMap["cpuUsage"].Error)

	_ = ctx.JSON(map[string]map[string]interface{}{
		"errAndReq": reqAndErrResp,
		"latency":   latencyResp,
		"cpuUsage":  cpuUsageResp,
	})

}

func (h *clusterHandler) GetPxData(ctx iris.Context) {
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParamDefault("st", "-10m")
	clusterIdx := ctx.URLParam("cluster_id")

	if err := validation.ValidateGetPxlData(clusterIdx, apiKey); err != nil {
		zkHttpResponse := utils.CreateErrorResponseWithStatusCode(err.Error)
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	pxResp := s.GetPxlData(ctx, clusterIdx, st, apiKey)
	utils.GenerateResponseAndReturn(ctx, pxResp)

}

// TODO: Refactor the resp and remove the method below
func getResp(resultStats *pxapi.ResultsStats, result interface{}, err *zkerrors.ZkError) map[string]interface{} {
	var x map[string]interface{}
	if result == nil || err != nil {
		x = map[string]interface{}{
			"results": nil,
			"stats":   nil,
			"status":  500,
		}
	} else {
		x = map[string]interface{}{
			"results": result,
			"stats":   resultStats,
			"status":  200,
		}
	}
	return x
}
