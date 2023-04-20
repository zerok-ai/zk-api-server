package handlerimplementation

import (
	"encoding/json"
	"fmt"
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

	Latency    *int    `json:"latency"`
	OtelFlag   *string `json:"otel_flag"`
	ReqBody    *string `json:"req_body"`
	ReqHeaders *string `json:"req_headers"`
	ReqMethod  *string `json:"req_method"`
	ReqPath    *string `json:"req_path"`
	RespBody   *string `json:"resp_body"`
	SpanId     *string `json:"span_id"`
	Time       *string `json:"time_"`
	TraceId    *string `json:"trace_id"`
	TraceState *string `json:"tracestate"`
	Type       *string `json:"type"`
}

func ConvertPixieDataToPixieTraceData(r *types.Record) PixieTraceData {
	var p = PixieTraceData{}

	p.Time, _ = utils.GetStringFromRecord("time_", r)
	p.Latency, _ = utils.GetIntegerFromRecord("latency", r)
	p.Type, _ = utils.GetStringFromRecord("type", r)
	p.TraceState, _ = utils.GetStringFromRecord("tracestate", r)
	p.TraceId, _ = utils.GetStringFromRecord("trace_id", r)
	p.SpanId, _ = utils.GetStringFromRecord("span_id", r)
	p.OtelFlag, _ = utils.GetStringFromRecord("otel_flag", r)
	p.ReqBody, _ = utils.GetStringFromRecord("req_body", r)
	p.RespBody, _ = utils.GetStringFromRecord("resp_body", r)
	p.ReqPath, _ = utils.GetStringFromRecord("req_path", r)
	p.ReqMethod, _ = utils.GetStringFromRecord("req_method", r)
	p.ReqHeaders, _ = utils.GetStringFromRecord("req_headers", r)

	s := Source{}
	d := Destination{}

	sStr, _ := utils.GetStringFromRecord("source", r)
	dStr, _ := utils.GetStringFromRecord("destination", r)

	json.Unmarshal([]byte(*sStr), &s)
	json.Unmarshal([]byte(*dStr), &d)

	p.Source = s
	p.Destination = d

	if !utils.IsEmpty(*p.ReqBody) && !utils.IsEmpty(*p.ReqHeaders) {
		sec := map[string]string{}
		if err := json.Unmarshal([]byte(*p.ReqHeaders), &sec); err != nil {
			fmt.Println("cannot covert req_headers to map, ", p.ReqHeaders)
		} else {
			if sec["Content-Encoding"] == "gzip" {
				v := utils.DecodeGzip(*p.ReqBody)
				p.ReqBody = &v
			}
		}
	}

	return p
}
