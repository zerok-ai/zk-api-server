package handlerimplementation

import (
	"encoding/json"
	"main/app/utils"
	"px.dev/pxapi/types"
)

type Service struct {
	ServiceName             string     `json:"service"`
	PodCount                *int       `json:"pod_count"`
	HttpLatencyIn           *Latencies `json:"http_latency_in"`
	HttpRequestThroughputIn *string    `json:"http_req_throughput_in"`
	HttpErrorRateIn         *float64   `json:"http_error_rate_in"`
	InboundConns            *float64   `json:"inbound_conns"`
	OutboundConns           *float64   `json:"outbound_conns"`
}

type Latencies struct {
	P01 float64 `json:"p01"`
	P10 float64 `json:"p10"`
	P25 float64 `json:"p25"`
	P50 float64 `json:"p50"`
	P75 float64 `json:"p75"`
	P90 float64 `json:"p90"`
	P99 float64 `json:"p99"`
}

func ConvertPixieDataToService(r *types.Record) Service {
	service := Service{}

	service.ServiceName = utils.GetStringFromRecord("service", r)
	service.HttpRequestThroughputIn = utils.GetStringPtrFromRecord("http_req_throughput_in", r)
	service.HttpLatencyIn = GetLatenciesPtr("http_latency_in", r)
	service.PodCount = utils.GetIntegerPtrFromRecord("pod_count", r)
	service.InboundConns = utils.GetFloat64PtrFromRecord("inbound_conns", r)
	service.OutboundConns = utils.GetFloat64PtrFromRecord("outbound_conns", r)
	service.HttpErrorRateIn = utils.GetFloat64PtrFromRecord("http_error_rate_in", r)

	return service
}

func GetLatencies(key string, r *types.Record) (Latencies, error) {
	v := utils.GetStringFromRecord(key, r)
	if v != "" {
		data := Latencies{}
		err := json.Unmarshal([]byte(v), &data)
		if err != nil {
			return Latencies{}, err
		}
		return data, nil
	}
	return Latencies{}, nil
}

func GetLatenciesPtr(key string, r *types.Record) *Latencies {
	v, err := GetLatencies(key, r)
	if err == nil {
		return &v
	}
	return &Latencies{}
}
