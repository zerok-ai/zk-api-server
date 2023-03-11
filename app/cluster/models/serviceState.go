package models

import (
	"main/app/cluster/utils"
	"px.dev/pxapi/types"
)

type ServiceStat struct {
	Time              string  `json:"time"`
	LatencyP50        int     `json:"latency_p50"`
	LatencyP90        int     `json:"latency_p90"`
	LatencyP99        int     `json:"latency_p99"`
	RequestThroughput float64 `json:"request_throughput"`
	ErrorRate         float64 `json:"error_rate"`
	ErrorsPerNs       float64 `json:"errors_per_ns"`
	BytesPerNs        float64 `json:"bytes_per_ns"`
}

func ConvertPixieDataToServiceStat(r *types.Record) ServiceStat {
	s := ServiceStat{}

	s.Time = r.GetDatum("time_").String()
	s.LatencyP50, _ = utils.GetInteger("latency_p50", r)
	s.LatencyP90, _ = utils.GetInteger("latency_p90", r)
	s.LatencyP99, _ = utils.GetInteger("latency_p99", r)
	s.RequestThroughput, _ = utils.GetFloat("request_throughput", r, 64)
	s.ErrorRate, _ = utils.GetFloat("error_rate", r, 64)
	s.ErrorsPerNs, _ = utils.GetFloat("errors_per_ns", r, 64)
	s.BytesPerNs, _ = utils.GetFloat("bytes_per_ns", r, 64)

	return s
}
