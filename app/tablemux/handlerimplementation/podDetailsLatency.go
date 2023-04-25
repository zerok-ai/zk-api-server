package handlerimplementation

type PodDetailsLatency struct {
	Time       *string `json:"time"`
	LatencyP50 *int    `json:"latency_p50"`
	LatencyP90 *int    `json:"latency_p90"`
	LatencyP99 *int    `json:"latency_p99"`
}
