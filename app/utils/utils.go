package utils

import (
	"fmt"
	"px.dev/pxapi/types"
	"regexp"
	"strconv"
	"strings"
)

var ResourceList = []string{"pod", "service", "workload", "namespace"}
var TimeUnitPxl = []string{"s", "m", "h", "d", "mon"}
var Actions = []string{"list", "map"}
var getNamespaceMethodTemplate = "get_namespace_data('%s')"
var getServiceMapMethodTemplate = "service_let_graph('%s')"
var getServiceListMethodTemplate = "my_fun('%s')"
var getPXDataMethodTemplate = "get_roi_data(\"%s\",%d,'%s')"
var getServiceDetailsMethodTemplate = "inbound_let_timeseries('%s', '%s')"

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
