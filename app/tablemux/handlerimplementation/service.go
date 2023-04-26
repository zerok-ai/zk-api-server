package handlerimplementation

type Service struct {
	ServiceName             *string    `json:"service"`
	PodCount                *int       `json:"pod_count"`
	HttpLatencyIn           *Latencies `json:"http_latency_in"`
	HttpRequestThroughputIn *float64   `json:"http_req_throughput_in"`
	HttpErrorRateIn         *float64   `json:"http_error_rate_in"`
	InboundConns            *float64   `json:"inbound_conns"`
	OutboundConns           *float64   `json:"outbound_conns"`
}

type Latencies struct {
	P01 *float64 `json:"p01"`
	P10 *float64 `json:"p10"`
	P25 *float64 `json:"p25"`
	P50 *float64 `json:"p50"`
	P75 *float64 `json:"p75"`
	P90 *float64 `json:"p90"`
	P99 *float64 `json:"p99"`
}
