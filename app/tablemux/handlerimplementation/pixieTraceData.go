package handlerimplementation

import (
	"encoding/json"
	"main/app/utils"
	"px.dev/pxapi/types"
)

type Source struct {
	Args   Args   `json:"args"`
	Label  string `json:"label"`
	Script string `json:"script"`
}

type Destination struct {
	Args   Args   `json:"args"`
	Label  string `json:"label"`
	Script string `json:"script"`
}

type Args struct {
	Pod       string `json:"pod"`
	StartTime string `json:"start_time"`
}

type PixieTraceData struct {
	Destination Destination `json:"destination"`
	Source      Source      `json:"source"`

	Latency    int    `json:"latency"`
	OtelFlag   string `json:"otel_flag"`
	ReqBody    string `json:"req_body"`
	ReqHeaders string `json:"req_headers"`
	ReqMethod  string `json:"req_method"`
	ReqPath    string `json:"req_path"`
	RespBody   string `json:"resp_body"`
	SpanId     string `json:"span_id"`
	Time       string `json:"time_"`
	TraceId    string `json:"trace_id"`
	TraceState string `json:"tracestate"`
	Type       string `json:"type"`
}

func ConvertPixieDataToPixieTraceData(r *types.Record) PixieTraceData {
	var p = PixieTraceData{}

	p.Time = utils.GetStringFromRecord("time_", r)
	p.Latency, _ = utils.GetIntegerFromRecord("latency", r)
	p.Type = utils.GetStringFromRecord("type", r)
	p.TraceState = utils.GetStringFromRecord("tracestate", r)
	p.TraceId = utils.GetStringFromRecord("trace_id", r)
	p.SpanId = utils.GetStringFromRecord("span_id", r)
	p.OtelFlag = utils.GetStringFromRecord("otel_flag", r)
	p.ReqBody = utils.GetStringFromRecord("req_body", r)
	p.RespBody = utils.GetStringFromRecord("resp_body", r)
	p.ReqPath = utils.GetStringFromRecord("req_path", r)
	p.ReqMethod = utils.GetStringFromRecord("req_method", r)
	p.ReqHeaders = utils.GetStringFromRecord("req_headers", r)

	s := Source{}
	d := Destination{}

	sStr := utils.GetStringFromRecord("source", r)
	dStr := utils.GetStringFromRecord("destination", r)

	json.Unmarshal([]byte(sStr), &s)
	json.Unmarshal([]byte(dStr), &d)

	p.Source = s
	p.Destination = d

	return p
}
