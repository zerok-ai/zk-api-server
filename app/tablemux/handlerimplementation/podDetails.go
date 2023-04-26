package handlerimplementation

type Status struct {
	Phase   *string `json:"phase"`
	Message *string `json:"message"`
	Reason  *string `json:"reason"`
	Ready   *bool   `json:"ready"`
}

type PodDetails struct {
	Pod        *string `json:"pod"`
	Service    *string `json:"service"`
	StartTime  *string `json:"start_time"`
	Containers *int    `json:"containers"`
	Status     Status  `json:"status"`
}
