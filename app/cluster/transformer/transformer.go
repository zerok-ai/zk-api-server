package transformer

import (
	"main/app/tablemux/handlerimplementation"
	"main/app/utils/zkerrors"
	"px.dev/pxapi"
)

type PixieHTTPResponse[T handlerimplementation.ItemType] struct {
	ResultStats *pxapi.ResultsStats `json:"stats"`
	Results     []T                 `json:"results"`
}

type PodDetailsPixieHTTPResponse struct {
	RequestAndError *PixieHTTPResponse[handlerimplementation.PodDetailsErrAndReq] `json:"requestAndError"`
	Latency         *PixieHTTPResponse[handlerimplementation.PodDetailsLatency]   `json:"latency"`
	CpuUsage        *PixieHTTPResponse[handlerimplementation.PodDetailsCpuUsage]  `json:"cpuUsage"`
}

func PixieResponseToHTTPResponse[T handlerimplementation.ItemType](results *pxapi.ScriptResults, mux *handlerimplementation.ItemMapMux[T], zkError *zkerrors.ZkError) *PixieHTTPResponse[T] {
	if zkError != nil {
		return nil
	}
	resp := PixieHTTPResponse[T]{}
	if results != nil {
		resp.ResultStats = results.Stats()
	}
	if mux != nil && mux.Table.Values != nil {
		resp.Results = mux.Table.Values
	}
	return &resp
}

func PixieResponseToPodDetailsHTTPResponse(reqAndErr *PixieHTTPResponse[handlerimplementation.PodDetailsErrAndReq], latency *PixieHTTPResponse[handlerimplementation.PodDetailsLatency], cpuUsage *PixieHTTPResponse[handlerimplementation.PodDetailsCpuUsage]) *PodDetailsPixieHTTPResponse {
	resp := PodDetailsPixieHTTPResponse{}

	resp.RequestAndError = reqAndErr
	resp.Latency = latency
	resp.CpuUsage = cpuUsage
	return &resp
}
