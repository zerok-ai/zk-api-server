package handlerimplementation

type ClusterDetailsMetaDataResponse struct {
	Status int `json:"status"`
	Data   struct {
		ApiKey struct {
			Typename string `json:"__typename"`
			Id       string `json:"id"`
			Key      string `json:"key"`
		} `json:"apiKey"`
		Clusters []ClusterDetailsFromResponse `json:"clusters"`
	} `json:"data"`
}

type ClusterDetailsFromResponse struct {
	Typename             string  `json:"__typename"`
	ClusterName          string  `json:"clusterName"`
	Id                   string  `json:"id"`
	LastHeartbeatMs      float64 `json:"lastHeartbeatMs"`
	NumInstrumentedNodes int     `json:"numInstrumentedNodes"`
	NumNodes             int     `json:"numNodes"`
	PrettyClusterName    string  `json:"prettyClusterName"`
	Status               string  `json:"status"`
	VizierVersion        string  `json:"vizierVersion"`
}
