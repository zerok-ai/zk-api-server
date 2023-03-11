package models

import (
	"main/app/cluster/utils"
	"px.dev/pxapi/types"
)

type ServiceMap struct {
	ResponderPod       string  `json:"responder_pod"`
	RequesterPod       string  `json:"requester_pod"`
	ResponderService   string  `json:"responder_service"`
	RequesterService   string  `json:"requester_service"`
	ResponderIP        string  `json:"responder_ip"`
	RequesterIP        string  `json:"requester_ip"`
	LatencyP50         float64 `json:"latency_p50"`
	LatencyP90         float64 `json:"latency_p90"`
	LatencyP99         float64 `json:"latency_p99"`
	RequestThroughput  float64 `json:"request_throughput"`
	ErrorRate          float64 `json:"error_rate"`
	InboundThroughput  float64 `json:"inbound_throughput"`
	OutboundThroughput float64 `json:"outbound_throughput"`
	ThroughputTotal    float64 `json:"throughput_total"`
}

func ConvertPixieDataToServiceMap(r *types.Record) ServiceMap {
	serviceMap := ServiceMap{}

	serviceMap.ResponderPod = utils.GetString("responder_pod", r)
	serviceMap.RequesterPod = utils.GetString("requestor_pod", r)
	serviceMap.ResponderService = utils.GetString("responder_service", r)
	serviceMap.RequesterService = utils.GetString("requestor_service", r)
	serviceMap.ResponderIP = utils.GetString("responder_ip", r)
	serviceMap.RequesterIP = utils.GetString("requestor_ip", r)
	serviceMap.LatencyP50, _ = utils.GetFloat("latency_p50", r, 64)
	serviceMap.LatencyP90, _ = utils.GetFloat("latency_p90", r, 64)
	serviceMap.LatencyP99, _ = utils.GetFloat("latency_p99", r, 64)
	serviceMap.RequestThroughput, _ = utils.GetFloat("request_throughput", r, 64)
	serviceMap.ErrorRate, _ = utils.GetFloat("error_rate", r, 64)
	serviceMap.InboundThroughput, _ = utils.GetFloat("inbound_throughput", r, 64)
	serviceMap.OutboundThroughput, _ = utils.GetFloat("outbound_throughput", r, 64)
	serviceMap.ThroughputTotal, _ = utils.GetFloat("throughput_total", r, 64)

	return serviceMap
}
