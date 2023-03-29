package handlerimplementation

import (
	"main/app/utils"
	"px.dev/pxapi/types"
)

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

func ConvertPixieDataToPodDetailsCpuUsage(r *types.Record) PodDetailsCpuUsage {
	var p = PodDetailsCpuUsage{}

	p.Time = utils.GetStringFromRecord("time_", r)
	p.CpuUsage = utils.GetStringFromRecord("cpu_usage", r)
	p.VSize = utils.GetStringFromRecord("vsize", r)
	p.Rss = utils.GetStringFromRecord("rss", r)
	p.TotalDiskWriteThroughput = utils.GetStringFromRecord("total_disk_write_throughput", r)
	p.TotalDiskReadThroughput = utils.GetStringFromRecord("total_disk_read_throughput", r)
	p.ActualDiskWriteThroughput = utils.GetStringFromRecord("actual_disk_write_throughput", r)
	p.ActualDiskReadThroughput = utils.GetStringFromRecord("actual_disk_read_throughput", r)
	p.Container = utils.GetStringFromRecord("container", r)

	return p
}
