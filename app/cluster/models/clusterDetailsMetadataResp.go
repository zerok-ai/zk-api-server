package models

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
