package handlerimplementation

type Destination struct {
	Args   Args    `json:"args"`
	Label  *string `json:"label"`
	Script *string `json:"script"`
}

type Args struct {
	Pod       *string `json:"pod"`
	StartTime string  `json:"start_time"`
}

type PixieTraceData struct {
	Destination Destination `json:"destination"`
	Source      Destination `json:"source"`

	Latency    *int    `json:"latency"`
	OtelFlag   *string `json:"otel_flag"`
	ReqBody    *string `json:"req_body"`
	ReqHeaders *string `json:"req_headers"`
	ReqMethod  *string `json:"req_method"`
	ReqPath    *string `json:"req_path"`
	RespBody   *string `json:"resp_body"`
	SpanId     *string `json:"span_id"`
	Time       *string `json:"time"`
	TraceId    *string `json:"trace_id"`
	TraceState *string `json:"tracestate"`
	Type       *string `json:"type"`
}
