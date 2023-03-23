package handlerimplementation

import (
	"encoding/json"
	"main/app/utils"
	"px.dev/pxapi/types"
)

type Status struct {
	Phase   string `json:"phase"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Ready   bool   `json:"ready"`
}

type PodDetails struct {
	Pod        string `json:"pod"`
	Service    string `json:"service"`
	StartTime  string `json:"startTime"`
	Containers *int   `json:"containers"`
	Status     Status `json:"status"`
}

func ConvertPixieDataToPodDetails(r *types.Record) PodDetails {
	var p = PodDetails{}

	p.Containers = utils.GetIntegerPtrFromRecord("containers", r)
	p.Pod = utils.GetStringFromRecord("pod", r)
	p.Service = utils.GetStringFromRecord("service", r)
	p.Pod = utils.GetStringFromRecord("pod", r)
	p.StartTime = utils.GetStringFromRecord("start_time", r)

	var s Status
	statusStr := utils.GetStringFromRecord("status", r)
	json.Unmarshal([]byte(statusStr), &s)
	p.Status = s

	return p
}
