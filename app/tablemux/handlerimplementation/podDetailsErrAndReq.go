package handlerimplementation

import (
	"main/app/utils"
	"px.dev/pxapi/types"
)

type PodDetailsErrAndReq struct {
	Time              *string  `json:"time"`
	Container         *string  `json:"container"`
	RequestThroughput *string  `json:"request_throughput"`
	ErrorsPerNs       *float64 `json:"errors_per_ns"`
	ErrorRate         *float64 `json:"error_rate"`
}

func ConvertPixieDataToPodDetailsErrAndReq(r *types.Record) PodDetailsErrAndReq {
	var p = PodDetailsErrAndReq{}

	p.Time, _ = utils.GetStringFromRecord("time_", r)
	p.Container, _ = utils.GetStringFromRecord("container", r)
	p.RequestThroughput, _ = utils.GetStringFromRecord("request_throughput", r)
	p.ErrorsPerNs, _ = utils.GetFloatFromRecord("errors_per_ns", r, 64)
	p.ErrorRate, _ = utils.GetFloatFromRecord("error_rate", r, 64)

	return p
}
