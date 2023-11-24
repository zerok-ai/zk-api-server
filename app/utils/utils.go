package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"px.dev/pxapi/types"

	"github.com/kataras/iris/v12/x/errors"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
)

var ResourceList = []string{"pod", "service", "workload", "namespace"}
var TimeUnitPxl = []string{"s", "m", "h", "d", "mon"}
var Actions = []string{"list", "map"}
var getNamespaceMethodTemplate = "get_namespace_data('%s')"
var getServiceMapMethodTemplate = "service_let_graph('%s')"
var getServiceListMethodTemplate = "my_fun('%s')"
var getHttpServiceListMethodTemplate = "http_svc('%s')"
var getMysqlServiceListMethodTemplate = "mysql_svc('%s')"
var getPgsqlServiceListMethodTemplate = "pgsql_svc('%s')"
var getPXDataMethodTemplate = "get_roi_data(\"%s\",%d,'%s')"
var getServiceDetailsMethodTemplate = "inbound_let_timeseries('%s', '%s')"
var getPodDetailsMethodTemplate = "pods('%s', '%s', '%s')"
var getPodDetailsForHTTPDataAndErrTemplate = "pod_details_inbound_request_timeseries_by_container('%s', '%s')"
var getPodDetailsForHTTPLatencyTemplate = "pod_details_inbound_latency_timeseries('%s', '%s')"
var getPodDetailsForCpuUsageTemplate = "pod_details_resource_timeseries('%s', '%s')"

const (
	ClusterIdHeader         = "Cluster-Id"
	HttpUtilsZkApiKeyHeader = "Zk-Api-Key"

	ClusterIdxPathParam     = "clusterIdx"
	IntegrationIdxPathParam = "integrationId"
	ObfuscationIdxPathParam = "obfuscationId"
	ScenarioIdxPathParam    = "scenarioIdx"

	LastSyncTS  = "last_sync_ts"
	Offset      = "offset"
	Limit       = "limit"
	Deleted     = "deleted"
	StartTime   = "st"
	Url         = "url"
	Name        = "name"
	Namespace   = "ns"
	ServiceName = "service_name"
	PodName     = "pod_name"
	ClusterId   = "cluster_id"
	Protocol    = "protocol"
	Version     = "version"
	File        = "file"

	EBPF = "EBPF"
	OTEL = "OTEL"

	Enable  = "enable"
	Disable = "disable"

	StatusError     = "error"
	ResponsePayload = "payload"
)

func GetDataByIdx(tag string, datatypeName string, r *types.Record) interface{} {
	var retVal any = nil
	var strRetVal, _ = GetStringFromRecord(tag, r)
	switch datatypeName {
	case "STRING":
		retVal, _ = GetStringFromRecord(tag, r)
	case "TIME64NS":
		retVal, _ = GetTimestampFromRecord(tag, r)
	case "BOOLEAN":
		retVal, _ = GetBooleanFromRecord(tag, r)
	case "INT64", "UINT128":
		retVal, _ = GetIntegerFromRecord(tag, r)
	case "FLOAT64":
		retVal, _ = GetFloatFromRecord(tag, r)
	case "DATA_TYPE_UNKNOWN":
		retVal, _ = GetStringFromRecord(tag, r)
	}

	jsonRetVal := map[string]interface{}{}

	err := json.Unmarshal([]byte(*strRetVal), &jsonRetVal)
	if err != nil {
		println(err)
		return retVal
	}

	return jsonRetVal

}

func GetStringFromRecord(key string, r *types.Record) (*string, error) {
	v := r.GetDatum(key)
	if v == nil {
		return nil, errors.New(fmt.Sprintf("key %s not found", key))
	}
	return zkCommon.ToPtr[string](v.String()), nil
}

func GetFloatFromRecord(key string, r *types.Record) (*float64, error) {
	dCasted := r.GetDatum(key).(*types.Float64Value)
	floatVal := zkCommon.Round(dCasted.Value(), 12)
	return &floatVal, nil
}

func GetIntegerFromRecord(key string, r *types.Record) (*int, error) {
	s, e := GetStringFromRecord(key, r)
	if s == nil {
		return nil, e
	}
	i, e := zkCommon.GetIntegerFromString(*s)
	return zkCommon.ToPtr[int](i), nil
}

func GetBooleanFromRecord(key string, r *types.Record) (*bool, error) {
	s, e := GetStringFromRecord(key, r)
	if s == nil || e != nil {
		return nil, e
	}
	boolValue, e := strconv.ParseBool(*s)
	return zkCommon.ToPtr[bool](boolValue), e
}

func GetTimestampFromRecord(key string, r *types.Record) (*string, error) {
	t, e := GetStringFromRecord(key, r)
	if t == nil || e != nil {
		return nil, e
	}
	strValue := string(*t)
	return zkCommon.ToPtr[string](strValue), nil
}

func GetFloatFromString(k string, b int) (float64, error) {
	return strconv.ParseFloat(k, b)
}

func GetNamespaceMethodSignature(st string) string {
	return fmt.Sprintf(getNamespaceMethodTemplate, st)
}

func GetServiceMapMethodSignature(st string) string {
	return fmt.Sprintf(getServiceMapMethodTemplate, st)
}

func GetServiceListMethodSignature(st string) string {
	return fmt.Sprintf(getServiceListMethodTemplate, st)
}

func GetHttpServiceListMethodSignature(st string) string {
	return fmt.Sprintf(getHttpServiceListMethodTemplate, st)
}

func GetMysqlServiceListMethodSignature(st string) string {
	return fmt.Sprintf(getMysqlServiceListMethodTemplate, st)
}

func GetPgsqlServiceListMethodSignature(st string) string {
	return fmt.Sprintf(getPgsqlServiceListMethodTemplate, st)
}

func GetPXDataSignature(head int, st, filter string) string {
	return fmt.Sprintf(getPXDataMethodTemplate, st, head, filter)
}

func GetServiceDetailsMethodSignature(st, serviceNameWithNs string) string {
	return fmt.Sprintf(getServiceDetailsMethodTemplate, st, serviceNameWithNs)
}

func GetPodDetailsMethodSignature(st, ns, serviceNameWithNs string) string {
	return fmt.Sprintf(getPodDetailsMethodTemplate, st, ns, serviceNameWithNs)
}

func GetPodDetailsForHTTPDataAndErrMethodSignature(st, podNameWithNs string) string {
	return fmt.Sprintf(getPodDetailsForHTTPDataAndErrTemplate, st, podNameWithNs)
}

func GetPodDetailsForHTTPLatencyMethodSignature(st, podNameWithNs string) string {
	return fmt.Sprintf(getPodDetailsForHTTPLatencyTemplate, st, podNameWithNs)
}

func GetPodDetailsForCpuUsageMethodSignature(st, podNameWithNs string) string {
	return fmt.Sprintf(getPodDetailsForCpuUsageTemplate, st, podNameWithNs)
}

func IsValidPxlTime(s string) bool {
	re := regexp.MustCompile("[0-9]+")
	d := re.FindAllString(s, -1)
	if len(d) != 1 {
		return false
	}

	t := strings.Split(s, d[0])
	var params = make([]string, 0)
	for _, v := range t {
		if !zkCommon.IsEmpty(v) {
			params = append(params, v)
		}
	}
	if len(params) == 2 {
		if !zkCommon.Contains(TimeUnitPxl, params[1]) || params[0] != "-" {
			return false
		}
	} else if len(params) == 1 {
		if !zkCommon.Contains(TimeUnitPxl, params[0]) {
			return false
		}
	} else {
		return false
	}

	return true
}

func ContainsValue(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
