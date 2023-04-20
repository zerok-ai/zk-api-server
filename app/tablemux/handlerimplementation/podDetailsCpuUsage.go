package handlerimplementation

import (
	"main/app/utils"
	"px.dev/pxapi/types"
)

type PodDetailsCpuUsage struct {
	Time                      *string `json:"time"`
	CpuUsage                  *string `json:"cpu_usage"`
	VSize                     *string `json:"vsize"`
	Rss                       *string `json:"rss"`
	TotalDiskWriteThroughput  *string `json:"total_disk_write_throughput"`
	TotalDiskReadThroughput   *string `json:"total_disk_read_throughput"`
	ActualDiskWriteThroughput *string `json:"actual_disk_write_throughput"`
	ActualDiskReadThroughput  *string `json:"actual_disk_read_throughput"`
	Container                 *string `json:"container"`
}

func ConvertPixieDataToPodDetailsCpuUsage(r *types.Record) PodDetailsCpuUsage {
	var p = PodDetailsCpuUsage{}

	p.Time, _ = utils.GetStringFromRecord("time_", r)
	p.CpuUsage, _ = utils.GetStringFromRecord("cpu_usage", r)
	p.VSize, _ = utils.GetStringFromRecord("vsize", r)
	p.Rss, _ = utils.GetStringFromRecord("rss", r)
	p.TotalDiskWriteThroughput, _ = utils.GetStringFromRecord("total_disk_write_throughput", r)
	p.TotalDiskReadThroughput, _ = utils.GetStringFromRecord("total_disk_read_throughput", r)
	p.ActualDiskWriteThroughput, _ = utils.GetStringFromRecord("actual_disk_write_throughput", r)
	p.ActualDiskReadThroughput, _ = utils.GetStringFromRecord("actual_disk_read_throughput", r)
	p.Container, _ = utils.GetStringFromRecord("container", r)

	return p
}
