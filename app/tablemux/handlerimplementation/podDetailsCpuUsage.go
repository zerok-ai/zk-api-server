package handlerimplementation

type PodDetailsCpuUsage struct {
	Time                      *string  `json:"time_"`
	CpuUsage                  *float64 `json:"cpu_usage"`
	VSize                     *float64 `json:"vsize"`
	Rss                       *float64 `json:"rss"`
	TotalDiskWriteThroughput  *float64 `json:"total_disk_write_throughput"`
	TotalDiskReadThroughput   *float64 `json:"total_disk_read_throughput"`
	ActualDiskWriteThroughput *float64 `json:"actual_disk_write_throughput"`
	ActualDiskReadThroughput  *float64 `json:"actual_disk_read_throughput"`
	Container                 *string  `json:"container"`
}
