package handlerimplementation

import (
	"encoding/json"
	"main/app/utils"
	"px.dev/pxapi/types"
)

type Status struct {
	Phase   *string `json:"phase"`
	Message *string `json:"message"`
	Reason  *string `json:"reason"`
	Ready   *bool   `json:"ready"`
}

type PodDetails struct {
	Pod        *string `json:"pod"`
	Service    *string `json:"service"`
	StartTime  *string `json:"startTime"`
	Containers *int    `json:"containers"`
	Status     Status  `json:"status"`
}

func ConvertPixieDataToPodDetails(r *types.Record) PodDetails {
	var p = PodDetails{}

	p.Containers, _ = utils.GetIntegerFromRecord("containers", r)
	p.Pod, _ = utils.GetStringFromRecord("pod", r)
	p.Service, _ = utils.GetStringFromRecord("service", r)
	p.Pod, _ = utils.GetStringFromRecord("pod", r)
	p.StartTime, _ = utils.GetStringFromRecord("start_time", r)

	var s Status
	statusStr, _ := utils.GetStringFromRecord("status", r)
	json.Unmarshal([]byte(*statusStr), &s)
	p.Status = s

	return p
}
