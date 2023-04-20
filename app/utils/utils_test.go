package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"px.dev/pxapi/types"
	"strconv"
	"testing"
)

func TestContains(t *testing.T) {
	strSlice := []string{"apple", "banana", "orange"}
	result := Contains(strSlice, "banana")
	assert.True(t, result)

	intSlice := []int{1, 2, 3}
	result = Contains(intSlice, 4)
	assert.False(t, result)

	emptySlice := []float64{}
	result = Contains(emptySlice, 5.5)
	assert.False(t, result)
}

func TestIsEmpty(t *testing.T) {
	assert.True(t, IsEmpty(""))

	assert.False(t, IsEmpty("hello"))
}

func TestGetIntegerFromString(t *testing.T) {
	inputStr := "1234"
	expected := 1234
	actual, err := GetIntegerFromString(inputStr)
	assert.Nil(t, err, "Unexpected error for valid input string")
	assert.Equal(t, expected, actual, "Unexpected result for valid input string")

	inputStr = "hello"
	expectedErr := strconv.NumError{
		Func: "Atoi",
		Num:  "hello",
		Err:  errors.New("invalid syntax"),
	}
	actual, err = GetIntegerFromString(inputStr)
	assert.Equal(t, &expectedErr, err, "Unexpected error for invalid input string")
	assert.Zero(t, actual, "Expected result to be 0 for invalid input string")
}

func TestGetFloatFromString(t *testing.T) {
	inputStr := "3.14159"
	inputBase := 64
	expected := 3.14159
	actual, err := GetFloatFromString(inputStr, inputBase)
	assert.Nil(t, err, "Unexpected error for valid input string")
	assert.Equal(t, expected, actual, "Unexpected result for valid input string")

	inputStr = "hello"
	expectedErr := strconv.NumError{
		Func: "ParseFloat",
		Num:  "hello",
		Err:  errors.New("invalid syntax"),
	}
	actual, err = GetFloatFromString(inputStr, inputBase)
	assert.Equal(t, &expectedErr, err, "Unexpected error for invalid input string")
	assert.Zero(t, actual, "Expected result to be 0 for invalid input string")
}

func TestStringToPtr(t *testing.T) {
	str := "test string"
	ptr := StringToPtr(str)
	assert.NotNil(t, ptr)
	assert.Equal(t, str, *ptr)
}

func TestIntToPtr(t *testing.T) {
	input := 1234
	expected := &input
	actual := IntToPtr(input)
	assert.Equal(t, *expected, *actual)

	input = -5678
	expected = &input
	actual = IntToPtr(input)
	assert.Equal(t, *expected, *actual)

	input = 0
	expected = &input
	actual = IntToPtr(input)
	assert.Equal(t, *expected, *actual)
}

func TestFloatToPtr(t *testing.T) {
	input := 1.2345
	expected := &input
	actual := FloatToPtr(input)

	assert.Equal(t, *expected, *actual)

	input = -3.1415
	expected = &input
	actual = FloatToPtr(input)
	assert.Equal(t, *expected, *actual)

	input = 0.0
	expected = &input
	actual = FloatToPtr(input)
	assert.Equal(t, *expected, *actual)
}
func TestGetNamespaceMethodSignature(t *testing.T) {
	st := "-10min"
	expectedResult := "get_namespace_data('-10min')"

	result := GetNamespaceMethodSignature(st)

	assert.Equal(t, expectedResult, result, "Method signature should match expected result")
}

func TestGetServiceMapMethodSignature(t *testing.T) {
	st := "-10min"
	expectedResult := "service_let_graph('-10min')"

	result := GetServiceMapMethodSignature(st)

	assert.Equal(t, expectedResult, result, "Method signature should match expected result")
}

func TestGetServiceListMethodSignature(t *testing.T) {
	st := "-10min"
	expectedResult := "my_fun('-10min')"

	result := GetServiceListMethodSignature(st)

	assert.Equal(t, expectedResult, result, "Method signature should match expected result")
}
func TestGetPXDataSignature(t *testing.T) {
	st := "-10min"
	head := 10
	filter := "{}"
	expectedResult := "get_roi_data(\"-10min\",10,'{}')"

	result := GetPXDataSignature(head, st, filter)

	assert.Equal(t, expectedResult, result, "Method signature should match expected result")

	filter = "{WHATEVER YOU GIVE HERE WILL BE USED AS FILTER}"
	expectedResult = "get_roi_data(\"-10min\",10,'{WHATEVER YOU GIVE HERE WILL BE USED AS FILTER}')"

	result = GetPXDataSignature(head, st, filter)

	assert.Equal(t, expectedResult, result, "Method signature should match expected result")
}

func TestGetServiceDetailsMethodSignature(t *testing.T) {
	st := "-10min"
	podNameWithNs := "namespace/pod-name"
	expectedResult := "inbound_let_timeseries('-10min', 'namespace/pod-name')"

	result := GetServiceDetailsMethodSignature(st, podNameWithNs)

	assert.Equal(t, expectedResult, result, "Method signature should match expected result")
}

func TestGetPodDetailsMethodSignature(t *testing.T) {
	st := "-10min"
	ns := "namespace"
	podName := "pod-name"
	expectedResult := "pods('-10min', 'namespace', 'pod-name')"

	result := GetPodDetailsMethodSignature(st, ns, podName)

	assert.Equal(t, expectedResult, result, "Method signature should match expected result")
}

func TestGetPodDetailsForHTTPDataAndErrMethodSignature(t *testing.T) {
	st := "-10min"
	podNameWithNs := "namespace/pod-name"
	expectedResult := "pod_details_inbound_request_timeseries_by_container('-10min', 'namespace/pod-name')"

	result := GetPodDetailsForHTTPDataAndErrMethodSignature(st, podNameWithNs)

	assert.Equal(t, expectedResult, result, "Method signature should match expected result")
}

func TestGetPodDetailsForHTTPLatencyMethodSignature(t *testing.T) {
	st := "-10min"
	podNameWithNs := "namespace/pod-name"
	expectedResult := "pod_details_inbound_latency_timeseries('-10min', 'namespace/pod-name')"

	result := GetPodDetailsForHTTPLatencyMethodSignature(st, podNameWithNs)

	assert.Equal(t, expectedResult, result, "Method signature should match expected result")
}

func TestGetPodDetailsForCpuUsageMethodSignature(t *testing.T) {
	st := "-10min"
	podNameWithNs := "namespace/pod-name"
	expectedResult := "pod_details_resource_timeseries('-10min', 'namespace/pod-name')"

	result := GetPodDetailsForCpuUsageMethodSignature(st, podNameWithNs)

	assert.Equal(t, expectedResult, result, "Method signature should match expected result")
}

func TestIsValidPxlTime(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"123s", true},
		{"5m", true},
		{"-5h", true},
		{"1d", true},
		{"30mon", true},
		{"30month", false},
		{"30mth", false},
		{"30", false},
		{"30.5m", false},
		{"35m5s", false},
		{"+5s", false},
		{"-5xyz", false},
	}

	for _, tc := range testCases {
		result := IsValidPxlTime(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestGetStringFromRecord(t *testing.T) {
	cs1 := types.ColSchema{
		Name:         "responder_pod",
		Type:         5,
		SemanticType: 400,
	}
	s1 := types.NewStringValue(&cs1)
	s1.ScanString("the_value_at_column_1")

	cs2 := types.ColSchema{
		Name:         "requester_pod",
		Type:         5,
		SemanticType: 400,
	}
	s2 := types.NewStringValue(&cs2)
	s2.ScanString("the_value_at_column_2")

	d := []types.Datum{s1, s2}

	table := types.TableMetadata{
		Name:    "table_name",
		ColInfo: nil,
		ColIdxByName: map[string]int64{
			"requester_pod": 1,
			"responder_pod": 0,
		},
	}

	mockRecord := &types.Record{Data: d, TableMetadata: &table}

	expectedOutput1 := "the_value_at_column_1"
	actualOutput1, _ := GetStringFromRecord("responder_pod", mockRecord)
	assert.Equal(t, expectedOutput1, *actualOutput1)

	expectedOutput2 := "the_value_at_column_2"
	actualOutput2, _ := GetStringFromRecord("requester_pod", mockRecord)
	assert.Equal(t, expectedOutput2, *actualOutput2)

	actualOutput3, _ := GetStringFromRecord("unknown_col", mockRecord)
	assert.Nil(t, actualOutput3)
}
func TestGetIntegerFromRecord(t *testing.T) {
	cs1 := types.ColSchema{
		Name:         "node_count",
		Type:         2,
		SemanticType: 1,
	}
	s1 := types.NewInt64Value(&cs1)
	s1.ScanInt64(1)

	d := []types.Datum{s1}

	table := types.TableMetadata{
		Name:    "table_name",
		ColInfo: nil,
		ColIdxByName: map[string]int64{
			"node_count": 0,
		},
	}

	mockRecord := &types.Record{Data: d, TableMetadata: &table}

	expectedOutput1 := 1
	actualOutput1, _ := GetIntegerFromRecord("node_count", mockRecord)
	assert.Equal(t, expectedOutput1, *actualOutput1)

	actualOutput3, _ := GetIntegerFromRecord("unknown_col", mockRecord)
	assert.Nil(t, actualOutput3)
}
func TestGetFloatFromRecord(t *testing.T) {
	cs1 := types.ColSchema{
		Name:         "latency_p50",
		Type:         4,
		SemanticType: 901,
	}
	s1 := types.NewFloat64Value(&cs1)
	s1.ScanFloat64(3.14)

	cs2 := types.ColSchema{
		Name:         "inbound_throughput",
		Type:         4,
		SemanticType: 903,
	}
	s2 := types.NewFloat64Value(&cs2)
	s2.ScanFloat64(0.34)

	d := []types.Datum{s1, s2}

	table := types.TableMetadata{
		Name:    "table_name",
		ColInfo: nil,
		ColIdxByName: map[string]int64{
			"latency_p50":        0,
			"inbound_throughput": 1,
		},
	}

	mockRecord := &types.Record{Data: d, TableMetadata: &table}

	expectedOutput1 := 3.14
	actualOutput1, _ := GetFloatFromRecord("latency_p50", mockRecord, 34)
	assert.Equal(t, expectedOutput1, *actualOutput1)

	expectedOutput2 := 0.34
	actualOutput2, _ := GetFloatFromRecord("inbound_throughput", mockRecord, 34)
	assert.Equal(t, expectedOutput2, *actualOutput2)

	actualOutput3, _ := GetFloatFromRecord("unknown_col", mockRecord, 34)
	assert.Nil(t, actualOutput3)
}

//func TestDecodeGzip(t *testing.T) {
//	input := "eJyr5lIAAqXU3MTMHCUrBaXElNzMPCOHxKRkveT8XCUdiHRBYnFxeX5RCkhFgRJXLQCVqg8A"
//	expectedOutput := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras ut risus quis mi rhoncus auctor. Nullam vehicula rutrum felis, quis suscipit odio interdum in. Fusce volutpat, mauris quis tempus malesuada, orci nibh mattis velit, vel aliquet odio enim vel justo. Etiam lacinia, mauris sed faucibus varius, tellus neque fringilla elit, ut placerat ipsum tellus id ante. Maecenas dapibus dignissim massa, eu pharetra eros hendrerit id. Fusce iaculis dui id est blandit luctus. Morbi eget ornare nunc. In hac habitasse platea dictumst.\n"
//	b := []byte(input)
//
//	actualOutput := readGzip(b)
//	assert.Equal(t, expectedOutput, actualOutput)
//
//	// Test the DecodeGzip function with a non-gzip input
//	//nonGzipInput := "This is not a gzip string"
//	//actualOutput = DecodeGzip(nonGzipInput)
//	//assert.Equal(t, nonGzipInput, actualOutput)
//}
