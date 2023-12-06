package utils

import (
	"strconv"
)

const (
	ClusterIdHeader         = "Cluster-Id"
	HttpUtilsZkApiKeyHeader = "Zk-Api-Key"

	ClusterIdxPathParam     = "clusterIdx"
	IntegrationIdxPathParam = "integrationId"
	ObfuscationIdxPathParam = "obfuscationId"
	ScenarioIdxPathParam    = "scenarioIdx"

	LastSyncTS  = "last_sync_ts"
	Offset      = "offset"
	Limit       = "limit"
	Deleted     = "deleted"
	StartTime   = "st"
	Url         = "url"
	Name        = "name"
	Namespace   = "ns"
	ServiceName = "service_name"
	ClusterId   = "cluster_id"
	Protocol    = "protocol"
	Version     = "version"
	File        = "file"

	EBPF = "EBPF"
	OTEL = "OTEL"

	Enable  = "enable"
	Disable = "disable"

	StatusError     = "error"
	ResponsePayload = "payload"
)

func GetFloatFromString(k string, b int) (float64, error) {
	return strconv.ParseFloat(k, b)
}

func ContainsValue(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
