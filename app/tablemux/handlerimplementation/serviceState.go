package handlerimplementation

type ServiceStat struct {
	Time              string   `json:"time"`
	LatencyP50        *float64 `json:"latency_p50"`
	LatencyP90        *float64 `json:"latency_p90"`
	LatencyP99        *float64 `json:"latency_p99"`
	RequestThroughput *float64 `json:"request_throughput"`
	ErrorRate         *float64 `json:"error_rate"`
	ErrorsPerNs       *float64 `json:"errors_per_ns"`
	BytesPerNs        *float64 `json:"bytes_per_ns"`
}