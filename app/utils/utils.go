package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"px.dev/pxapi/types"
)

var ResourceList = []string{"pod", "service", "workload", "namespace"}
var TimeUnitPxl = []string{"s", "m", "h", "d", "mon"}
var Actions = []string{"list", "map"}
var getNamespaceMethodTemplate = "get_namespace_data('%s')"
var getServiceMapMethodTemplate = "service_let_graph('%s')"
var getServiceListMethodTemplate = "my_fun('%s')"
var getPXDataMethodTemplate = "get_roi_data(\"%s\",%d,'%s')"
var getServiceDetailsMethodTemplate = "inbound_let_timeseries('%s', '%s')"
var getPodDetailsMethodTemplate = "pods('%s', '%s', '%s')"
var getPodDetailsForHTTPDataAndErrTemplate = "pod_details_inbound_request_timeseries_by_container('%s', '%s')"
var getPodDetailsForHTTPLatencyTemplate = "pod_details_inbound_latency_timeseries('%s', '%s')"
var getPodDetailsForCpuUsageTemplate = "pod_details_resource_timeseries('%s', '%s')"

func GetData(tag string, datatypeName string, r *types.Record) interface{} {
	var retVal any = nil
	switch datatypeName {
	case "STRING", "DATA_TYPE_UNKNOWN", "BOOLEAN", "TIME64NS":
		retVal = GetStringFromRecord(tag, r)
	case "INT64", "UINT128":
		retVal, _ = GetIntegerFromRecord(tag, r)
	case "FLOAT64":
		retVal, _ = GetFloatFromRecord(tag, r, 64)
	}
	return retVal
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func IsEmpty(v string) bool {
	return len(v) == 0
}

func StringToPtr(v string) *string {
	return &v
}

func GetStringFromRecord(key string, r *types.Record) string {
	return r.GetDatum(key).String()
}

func GetFloatFromRecord(key string, r *types.Record, bitSize int) (float64, error) {
	return GetFloatFromString(GetStringFromRecord(key, r), bitSize)
}

func GetIntegerFromRecord(key string, r *types.Record) (int, error) {
	return GetIntegerFromString(GetStringFromRecord(key, r))
}

func GetStringPtrFromRecord(key string, r *types.Record) *string {
	return StringToPtr(GetStringFromRecord(key, r))
}

func GetFloat64PtrFromRecord(key string, r *types.Record) *float64 {
	v, err := GetFloatFromRecord(key, r, 64)
	if err != nil {
		return nil
	} else {
		return &v
	}
}

func GetIntegerPtrFromRecord(key string, r *types.Record) *int {
	v, err := GetIntegerFromRecord(key, r)
	if err != nil {
		return nil
	} else {
		return &v
	}
}

func GetIntegerFromString(k string) (int, error) {
	return strconv.Atoi(k)
}

func GetFloatFromString(k string, b int) (float64, error) {
	return strconv.ParseFloat(k, b)
}

func IntToPtr(v int) *int {
	return &v
}

func FloatToPtr(v float64) *float64 {
	return &v
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
	if len(t) == 3 {
		if !Contains(TimeUnitPxl, t[2]) && t[0] != "-" {
			return false
		}
	} else if len(t) == 2 {
		if !Contains(TimeUnitPxl, t[1]) {
			return false
		}
	} else {
		return false
	}

	return true
}

func DecodeGzip(s string) string {
	b := []byte(s)
	reader := bytes.NewReader(b)
	r, err := gzip.NewReader(reader)
	defer func(r *gzip.Reader) {
		err := r.Close()
		if err != nil {

		}
	}(r)

	if err != nil {
		log.Printf("Error while decoding gzip string %s\n", s)
		log.Printf(err.Error())
	}

	output, err := io.ReadAll(r)
	if err != nil {
		log.Printf("Error while reading gzip string %s\n", s)
		log.Printf(err.Error())
	}

	return string(output)
}
