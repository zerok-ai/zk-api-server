package handlerimplementation

type PodDetailsCpuUsage struct {
	Time                      string `json:"time"`
	CpuUsage                  string `json:"cpu_usage"`
	VSize                     string `json:"vsize"`
	Rss                       string `json:"rss"`
	TotalDiskWriteThroughput  string `json:"total_disk_write_throughput"`
	TotalDiskReadThroughput   string `json:"total_disk_read_throughput"`
	ActualDiskWriteThroughput string `json:"actual_disk_write_throughput"`
	ActualDiskReadThroughput  string `json:"actual_disk_read_throughput"`
	Container                 string `json:"container"`
}
