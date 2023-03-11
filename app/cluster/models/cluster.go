package models

type Cluster struct {
	Nickname  string  `json:"nickname"`
	Domain    string  `json:"domain"`
	ApiKey    string  `json:"api_key"`
	ClusterId string  `json:"cluster_id"`
	Id        *string `json:"id,omitempty"`
}

var ClusterMap = map[string]Cluster{
	"1": {
		"devpx07",
		"devpx07.getanton.com",
		"px-api-13c7defc-7e35-4ce6-a9e4-f6a623137e7c",
		"8941bee5-6621-442c-9de8-309ecfb5ec22",
		nil,
	},
	"2": {
		"avinpx06",
		"avinpx06.getanton.com",
		"px-api-d8dcdc33-bacc-44f3-b92e-220599139609",
		"399c0e93-e34f-4b99-9a16-e9315147e8ba",
		nil,
	},
	"3": {
		"sockpx01",
		"sockpx01.getanton.com",
		"px-api-03c9d006-3490-416c-b53a-64044241c420",
		"65a8228d-48b0-4ef8-b174-b460f639df65",
		nil,
	},
}
