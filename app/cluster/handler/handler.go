package handler

import (
	"github.com/kataras/iris/v12"
	"main/app/cluster/service"
	"main/app/cluster/transformer"
	"main/app/cluster/validation"
	"main/app/tablemux"
	"main/app/tablemux/handlerimplementation"
	"main/app/utils"
)

type ClusterHandler interface {
	GetServiceDetailsList(ctx iris.Context)
	GetServiceDetailsMap(ctx iris.Context)
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

// GetServiceDetailsList Returns all services in the provided cluster with their details
//
//	@Summary		Get all services' details
//	@Description	Returns all services in the provided cluster with their details
//	@Tags			cluster data
//	@Produce		json
//	@Success		200 {object} utils.ZkHttpResponse[transformer.PixieHTTPResponse[handlerimplementation.Service]]
//	@Router			/u/cluster/{clusterIdx}/service/list [get]
func (h *clusterHandler) GetServiceDetailsList(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParam("st")

	if err := validation.ValidateGetResourceDetailsApi(st, apiKey); err != nil {
		zkHttpResponse := utils.ZkHttpResponseBuilder[any]{}.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	resp, zkError := s.GetServiceDetailsList(ctx, clusterIdx, st, apiKey)
	utils.SetResponseInCtxAndReturn[transformer.PixieHTTPResponse[handlerimplementation.Service]](ctx, resp, zkError)
}

// GetServiceDetailsMap Returns the data between two services that directly interacts with each other
//
//	@Summary		Get all services' map
//	@Description	Returns the data between two services that directly interacts with each other
//	@Tags			cluster data
//	@Produce		json
//	@Success		200 {object} utils.ZkHttpResponse[transformer.PixieHTTPResponse[handlerimplementation.ServiceMap]]
//	@Router			/u/cluster/{clusterIdx}/service/map [get]
func (h *clusterHandler) GetServiceDetailsMap(ctx iris.Context) {
	apiKey := ctx.GetHeader("ZK_API_KEY")
	clusterIdx := ctx.Params().Get("clusterIdx")
	st := ctx.URLParam("st")

	if err := validation.ValidateGetResourceDetailsApi(st, apiKey); err != nil {
		zkHttpResponse := utils.ZkHttpResponseBuilder[any]{}.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	resp, zkError := s.GetServiceDetailsMap(ctx, clusterIdx, st, apiKey)
	utils.SetResponseInCtxAndReturn[transformer.PixieHTTPResponse[handlerimplementation.ServiceMap]](ctx, resp, zkError)
}

// GetServiceDetails Returns the data of a service under a cluster
//
//	@Summary		Get all data from a service
//	@Description	Returns the data of a service
//	@Tags			cluster data
//	@Produce		json
//	@Success		200 {object} utils.ZkHttpResponse[transformer.PixieHTTPResponse[handlerimplementation.ServiceStat]]
//	@Router			/u/cluster/{clusterIdx}/service/details [get]
func (h *clusterHandler) GetServiceDetails(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	apiKey := ctx.GetHeader("ZK_API_KEY")
	serviceName := ctx.URLParam("name")
	ns := ctx.URLParam("ns")
	st := ctx.URLParam("st")

	if err := validation.ValidateGraphDetailsApi(serviceName, ns, st, apiKey); err != nil {
		zkHttpResponse := utils.ZkHttpResponseBuilder[any]{}.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	resp, zkError := s.GetServiceDetails(ctx, clusterIdx, serviceName, ns, st, apiKey)
	utils.SetResponseInCtxAndReturn[transformer.PixieHTTPResponse[handlerimplementation.ServiceStat]](ctx, resp, zkError)
}

// GetPodList Returns a list of all the pods under a cluster and service
//
//	@Summary		Get all pods under a service
//	@Description	Returns a list of all the pods under a cluster and service
//	@Tags			cluster data
//	@Produce		json
//	@Success		200 {object} utils.ZkHttpResponse[transformer.PixieHTTPResponse[handlerimplementation.PodDetails]]
//	@Router			/u/cluster/{clusterIdx}/pod/list [get]
func (h *clusterHandler) GetPodList(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParam("st")
	serviceName := ctx.URLParam("service_name")
	ns := ctx.URLParam("ns")

	if err := validation.ValidateGraphDetailsApi(serviceName, ns, st, apiKey); err != nil {
		zkHttpResponse := utils.ZkHttpResponseBuilder[any]{}.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	resp, zkError := s.GetPodList(ctx, clusterIdx, serviceName, ns, st, apiKey)
	utils.SetResponseInCtxAndReturn[transformer.PixieHTTPResponse[handlerimplementation.PodDetails]](ctx, resp, zkError)
}

// GetPodDetails Returns time-series data for the given pod
//
//	@Summary		Returns time-series data for the given pod
//	@Description	Returns time-series data for the given pod for Request And Error, latency and cpu usage
//	@Tags			cluster data
//	@Produce		json
//	@Success		200 {object} utils.ZkHttpResponse[transformer.PodDetailsPixieHTTPResponse]
//	@Router			/u/cluster/{clusterIdx}/pod/details [get]
func (h *clusterHandler) GetPodDetails(ctx iris.Context) {
	clusterIdx := ctx.Params().Get("clusterIdx")
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParam("st")
	podName := ctx.URLParam("pod_name")
	ns := ctx.URLParam("ns")

	if err := validation.ValidatePodDetailsApi(ctx, podName, ns, st, apiKey); err != nil {
		zkHttpResponse := utils.ZkHttpResponseBuilder[any]{}.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}
	resp, zkError := s.GetPodDetailsTimeSeries(ctx, clusterIdx, podName, ns, st, apiKey)
	utils.SetResponseInCtxAndReturn[transformer.PodDetailsPixieHTTPResponse](ctx, resp, zkError)
}

// GetPxData Returns pixie data for a given cluster
//
//	@Summary		Get pixie data
//	@Description	Returns pixie data for a given cluster
//	@Tags			cluster data
//	@Produce		json
//	@Success		200 {object} utils.ZkHttpResponse[transformer.PixieHTTPResponse[handlerimplementation.PixieTraceData]]
//	@Router			/u/cluster/traces [get]
func (h *clusterHandler) GetPxData(ctx iris.Context) {
	apiKey := ctx.GetHeader("ZK_API_KEY")
	st := ctx.URLParamDefault("st", "-10m")
	clusterIdx := ctx.URLParam("cluster_id")

	if err := validation.ValidateGetPxlData(clusterIdx, apiKey); err != nil {
		zkHttpResponse := utils.ZkHttpResponseBuilder[any]{}.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	resp, zkError := s.GetPxlData(ctx, clusterIdx, st, apiKey)
	utils.SetResponseInCtxAndReturn[transformer.PixieHTTPResponse[handlerimplementation.PixieTraceData]](ctx, resp, zkError)
}
