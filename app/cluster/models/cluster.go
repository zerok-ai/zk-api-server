package models

type ClusterDetails struct {
	Nickname    string  `json:"nickname,omitempty"`
	Domain      string  `json:"domain,omitempty"`
	ApiKey      string  `json:"api_key,omitempty"`
	ClusterId   string  `json:"cluster_id,omitempty"`
	Id          *string `json:"id,omitempty"`
	ClusterName string  `json:"cluster_name,omitempty"`
	Status      string  `json:"status,omitempty"`
}

var ClusterMap = map[string]ClusterDetails{}
