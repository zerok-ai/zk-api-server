package cluster

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
		"px-api-4ccea853-5f4e-4127-969d-5a4f288a47ac",
		"b60b7536-acad-44ad-96fe-604a528660ce",
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
