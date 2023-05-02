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

func (h *clusterHandler) GetResourceDetailsList(ctx iris.Context) {
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

func (h *clusterHandler) GetResourceDetailsMap(ctx iris.Context) {
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
