package handlerimplementation

type ServiceMap struct {
	ResponderPod       *string  `json:"responder_pod"`
	RequesterPod       *string  `json:"requester_pod"`
	ResponderService   *string  `json:"responder_service"`
	RequesterService   *string  `json:"requester_service"`
	ResponderIP        *string  `json:"responder_ip"`
	RequesterIP        *string  `json:"requester_ip"`
	LatencyP50         *float64 `json:"latency_p50"`
	LatencyP90         *float64 `json:"latency_p90"`
	LatencyP99         *float64 `json:"latency_p99"`
	RequestThroughput  *float64 `json:"request_throughput"`
	ErrorRate          *float64 `json:"error_rate"`
	InboundThroughput  *float64 `json:"inbound_throughput"`
	OutboundThroughput *float64 `json:"outbound_throughput"`
	ThroughputTotal    *float64 `json:"throughput_total"`
}
