package models

type Cluster struct {
	Nickname    string  `json:"nickname,omitempty"`
	Domain      string  `json:"domain,omitempty"`
	ApiKey      string  `json:"api_key,omitempty"`
	ClusterId   string  `json:"cluster_id,omitempty"`
	Id          *string `json:"id,omitempty"`
	ClusterName string  `json:"cluster_name,omitempty"`
	Status      string  `json:"status,omitempty"`
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

func (r *ClusterDetailsFromResponse) FromResponseToDomainClusterDetails() Cluster {
	return Cluster{
		Nickname:    r.PrettyClusterName,
		ClusterId:   r.Id,
		ClusterName: r.ClusterName,
		Status:      r.Status,
	}
}

var ClusterMap = map[string]Cluster{}

//func init() {
//
//	path, err := os.Getwd()
//	if err != nil {
//		log.Println(err)
//		os.Exit(2)
//	}
//	configFilePath := path + "/app/cluster/models/del.txt"
//	//configFilePath := "/opt/cluster.conf"
//
//	jsonFile, err := os.Open(configFilePath)
//
//	if err != nil {
//		log.Println(err)
//		os.Exit(2)
//	}
//	defer jsonFile.Close()
//
//	err = json.NewDecoder(jsonFile).Decode(&ClusterMap)
//	if err != nil {
//		log.Println(err)
//		os.Exit(2)
//	}
//	return
//}
