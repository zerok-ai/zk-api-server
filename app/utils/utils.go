package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12/x/errors"
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

var HTTP_DEBUG = false

func Init(httpDebug bool) {
	HTTP_DEBUG = httpDebug

}

func GetDataByIdx(idx int, datatypeName string, r *types.Record) interface{} {
	tag := r.TableMetadata.ColInfo[idx].Name
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
		retVal, _ = GetFloatFromRecord(idx, r, 64)
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

//func StringToPtr(v string) *string {
//	return &v
//}

func ToPtr[T any](arg T) *T {
	return &arg
}

func GetStringFromRecord(key string, r *types.Record) (*string, error) {
	v := r.GetDatum(key)
	if v == nil {
		return nil, errors.New(fmt.Sprintf("key %s not found", key))
	}
	return ToPtr[string](v.String()), nil
}

func GetFloatFromRecord(idx int, r *types.Record, bitSize int) (*float64, error) {
	dCasted, _ := r.Data[idx].(*types.Float64Value)
	var floatVal float64 = dCasted.Value()
	return &floatVal, nil
}

func GetIntegerFromRecord(key string, r *types.Record) (*int, error) {
	s, e := GetStringFromRecord(key, r)
	if s == nil {
		return nil, e
	}
	i, e := GetIntegerFromString(*s)
	return ToPtr[int](i), nil
}

func GetBooleanFromRecord(key string, r *types.Record) (*bool, error) {
	s, e := GetStringFromRecord(key, r)
	if s == nil || e != nil {
		return nil, e
	}
	boolValue, e := strconv.ParseBool(*s)
	return ToPtr[bool](boolValue), e
}

func GetTimestampFromRecord(key string, r *types.Record) (*string, error) {
	t, e := GetStringFromRecord(key, r)
	if t == nil || e != nil {
		return nil, e
	}
	strValue := string(*t)
	return ToPtr[string](strValue), nil
}

func GetIntegerFromString(k string) (int, error) {
	return strconv.Atoi(k)
}

func GetFloatFromString(k string, b int) (float64, error) {
	return strconv.ParseFloat(k, b)
}

//func IntToPtr(v int) *int {
//	return &v
//}
//
//func FloatToPtr(v float64) *float64 {
//	return &v
//}

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
	var params = make([]string, 0)
	for _, v := range t {
		if !IsEmpty(v) {
			params = append(params, v)
		}
	}
	if len(params) == 2 {
		if !Contains(TimeUnitPxl, params[1]) || params[0] != "-" {
			return false
		}
	} else if len(params) == 1 {
		if !Contains(TimeUnitPxl, params[0]) {
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
