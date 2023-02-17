package cluster

func ValidClusterId(clusterId string) bool {
	// must be present in ClusterMap.
	_, exist := ClusterMap[clusterId]
	return exist
}
