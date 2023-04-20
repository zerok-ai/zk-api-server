package handlerimplementation

import (
	"main/app/utils"
	"px.dev/pxapi/types"
)

type PodDetailsLatency struct {
	Time       *string `json:"time"`
	LatencyP50 *int    `json:"latency_p50"`
	LatencyP90 *int    `json:"latency_p90"`
	LatencyP99 *int    `json:"latency_p99"`
}

func ConvertPixieDataToPodDetailsLatency(r *types.Record) PodDetailsLatency {
	var p = PodDetailsLatency{}

	p.Time, _ = utils.GetStringFromRecord("time_", r)
	p.LatencyP50, _ = utils.GetIntegerFromRecord("latency_p50", r)
	p.LatencyP90, _ = utils.GetIntegerFromRecord("latency_p90", r)
	p.LatencyP99, _ = utils.GetIntegerFromRecord("latency_p99", r)

	return p
}
