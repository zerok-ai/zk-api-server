package models

import (
	"encoding/json"
	"main/app/cluster/utils"
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

	p.Time = utils.GetString("time_", r)
	p.Type = utils.GetString("type", r)
	p.TraceState = utils.GetString("tracestate", r)
	p.TraceId = utils.GetString("trace_id", r)
	p.SpanId = utils.GetString("span_id", r)
	p.OtelFlag = utils.GetString("otel_flag", r)
	p.ReqMethod = utils.GetString("req_body", r)
	p.RespBody = utils.GetString("resp_body", r)
	p.ReqPath = utils.GetString("req_path", r)
	p.ReqMethod = utils.GetString("req_method", r)
	p.ReqHeaders = utils.GetString("req_headers", r)

	s := Source{}
	d := Destination{}

	sStr := utils.GetString("source", r)
	dStr := utils.GetString("source", r)

	json.Unmarshal([]byte(sStr), &s)
	json.Unmarshal([]byte(dStr), &d)

	p.Source = s
	p.Destination = d

	return p
}
