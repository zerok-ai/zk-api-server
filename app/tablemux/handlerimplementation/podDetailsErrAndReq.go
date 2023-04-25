package handlerimplementation

type PodDetailsErrAndReq struct {
	Time              *string  `json:"time"`
	Container         *string  `json:"container"`
	RequestThroughput *string  `json:"request_throughput"`
	ErrorsPerNs       *float64 `json:"errors_per_ns"`
	ErrorRate         *float64 `json:"error_rate"`
}
