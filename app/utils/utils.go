package utils

import "fmt"

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
func GetServiceStatsMethodSignature(st, serviceNameWithNs string) string {
	return fmt.Sprintf(getServiceStatsMethodTemplate, st, serviceNameWithNs)
}

var ResourceList = []string{"pod", "service", "workload", "namespace"}
var Actions = []string{"list", "map"}
var getNamespaceMethodTemplate = "get_namespace_data('%s')"
var getServiceMapMethodTemplate = "service_let_graph('%s')"
var getServiceListMethodTemplate = "my_fun('%s')"
var getServiceStatsMethodTemplate = "inbound_let_timeseries('%s', '%s')"
