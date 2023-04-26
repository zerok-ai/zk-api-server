package handlerimplementation

type PodDetailsLatency struct {
	Time       *string  `json:"time_"`
	LatencyP50 *float64 `json:"latency_p50"`
	LatencyP90 *float64 `json:"latency_p90"`
	LatencyP99 *float64 `json:"latency_p99"`
}
