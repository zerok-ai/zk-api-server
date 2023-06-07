package utils

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	zkcommon "github.com/zerok-ai/zk-utils-go/common"
	"math"
	"px.dev/pxapi/proto/vizierpb"
	"px.dev/pxapi/types"
	"strconv"
	"testing"
)

func TestContains(t *testing.T) {
	strSlice := []string{"apple", "banana", "orange"}
	result := zkcommon.Contains(strSlice, "banana")
	assert.True(t, result)

	intSlice := []int{1, 2, 3}
	result = zkcommon.Contains(intSlice, 4)
	assert.False(t, result)

	emptySlice := []float64{}
	result = zkcommon.Contains(emptySlice, 5.5)
	assert.False(t, result)
}

func TestGetDataByIdx(t *testing.T) {
	// Create a test record
	cs0 := types.ColSchema{
		Name:         "zeroth_column_integer",
		Type:         2,
		SemanticType: 1,
	}
	s0 := types.NewInt64Value(&cs0)
	s0.ScanInt64(42)

	cs1 := types.ColSchema{
		Name:         "first_column_string",
		Type:         5,
		SemanticType: 300,
	}
	s1 := types.NewStringValue(&cs1)
	s1.ScanString("test string")

	cs2 := types.ColSchema{
		Name:         "second_column_bool",
		Type:         1,
		SemanticType: 1,
	}
	s2 := types.NewBooleanValue(&cs2)
	s2.ScanBool(true)

	cs3 := types.ColSchema{
		Name:         "third_column_float",
		Type:         4,
		SemanticType: 901,
	}
	s3 := types.NewFloat64Value(&cs3)
	s3.ScanFloat64(3.14159)

	cs4 := types.ColSchema{
		Name:         "fourth_column_time",
		Type:         6,
		SemanticType: 1,
	}
	s4 := types.NewTime64NSValue(&cs4)
	s4.ScanInt64(1683190341)

	cs5 := types.ColSchema{
		Name:         "fifth_column_datatype_unknown",
		Type:         0,
		SemanticType: 402,
	}
	s5 := types.NewStringValue(&cs5)
	s5.ScanString(`{"phase":"Running","message":"","reason":"","ready":true}`)

	tableMetaData := types.TableMetadata{
		Name:    "table_meta_data_1",
		ColInfo: []types.ColSchema{cs0, cs1, cs2, cs3, cs4, cs5},
		ColIdxByName: map[string]int64{
			"zeroth_column_integer":         0,
			"first_column_string":           1,
			"second_column_bool":            2,
			"third_column_float":            3,
			"fourth_column_time":            4,
			"fifth_column_datatype_unknown": 5,
		},
	}

	data := []types.Datum{s0, s1, s2, s3, s4, s5}
	// Test cases
	testCases := []struct {
		name     string
		expected interface{}
	}{
		{
			name:     "int64",
			expected: 42,
		},
		{
			name:     "string",
			expected: "test string",
		},
		{
			name:     "boolean",
			expected: true,
		},
		{
			name:     "float64",
			expected: 3.14159,
		},
		{
			name:     "time",
			expected: "1970-01-01 05:30:01.683190341 +0530 IST",
		},
		{
			name: "unknown data type",
			expected: map[string]interface{}{
				"phase":   "Running",
				"message": "",
				"reason":  "",
				"ready":   true,
			},
		},
	}

	r := types.Record{
		Data:          data,
		TableMetadata: &tableMetaData,
	}

	for i, tc := range testCases {
		dataType := vizierpb.DataType_name[int32(r.TableMetadata.ColInfo[i].Type)]
		result := GetDataByIdx(r.TableMetadata.ColInfo[i].Name, dataType, &r)

		switch dataType {
		case "STRING":
			assert.Equal(t, GetDataByIdx_HelperFunc[string](tc.expected), *GetDataByIdx_HelperFunc[*string](result))
		case "TIME64NS":
			assert.Equal(t, GetDataByIdx_HelperFunc[string](tc.expected), *GetDataByIdx_HelperFunc[*string](result))
		case "BOOLEAN":
			assert.Equal(t, GetDataByIdx_HelperFunc[bool](tc.expected), *GetDataByIdx_HelperFunc[*bool](result))
		case "INT64", "UINT128":
			assert.Equal(t, GetDataByIdx_HelperFunc[int](tc.expected), *GetDataByIdx_HelperFunc[*int](result))
		case "FLOAT64":
			assert.Equal(t, GetDataByIdx_HelperFunc[float64](tc.expected), *GetDataByIdx_HelperFunc[*float64](result))
		case "DATA_TYPE_UNKNOWN":
			assert.Equal(t, result.(map[string]interface{}), tc.expected.(map[string]interface{}))
		}

		// Check that the result can be marshalled to JSON
		_, err := json.Marshal(result)
		assert.NoError(t, err)
	}
}

func GetDataByIdx_HelperFunc[T any](val interface{}) T {
	x := val.(T)
	return x
}

func TestIsEmpty(t *testing.T) {
	assert.True(t, zkcommon.IsEmpty(""))

	assert.False(t, zkcommon.IsEmpty("hello"))
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
		Name:         "requestor_pod",
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
			"requestor_pod": 1,
			"responder_pod": 0,
		},
	}

	mockRecord := &types.Record{Data: d, TableMetadata: &table}

	expectedOutput1 := "the_value_at_column_1"
	actualOutput1, _ := GetStringFromRecord("responder_pod", mockRecord)
	assert.Equal(t, expectedOutput1, *actualOutput1)

	expectedOutput2 := "the_value_at_column_2"
	actualOutput2, _ := GetStringFromRecord("requestor_pod", mockRecord)
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

func TestRound(t *testing.T) {
	testCases := []struct {
		input     float64
		expected  float64
		precision int
	}{
		{2.345, 2.35, 2},
		{5.6789, 5.679, 3},
		{1.1111, 1, 0},
		{0, 0, 2},
		{math.Pi, 3.1416, 4},
	}

	for _, tc := range testCases {
		result := Round(tc.input, tc.precision)
		if result != tc.expected {
			t.Errorf("Expected Round(%f, %d) to be %f, but got %f", tc.input, tc.precision, tc.expected, result)
		}
	}
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
	actualOutput1, _ := GetFloatFromRecord("latency_p50", mockRecord)
	assert.Equal(t, expectedOutput1, *actualOutput1)

	expectedOutput2 := 0.34
	actualOutput2, _ := GetFloatFromRecord("inbound_throughput", mockRecord)
	assert.Equal(t, expectedOutput2, *actualOutput2)
}

//func TestDecodeGzip(t *testing.T) {
//	input := "eJyr5lIAAqXU3MTMHCUrBaXElNzMPCOHxKRkveT8XCUdiHRBYnFxeX5RCkhFgRJXLQCVqg8A"
//	expectedOutput := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras ut risus quis mi rhoncus auctor. Nullam vehicula rutrum felis, quis suscipit odio interdum in. Fusce volutpat, mauris quis tempus malesuada, orci nibh mattis velit, vel aliquet odio enim vel justo. Etiam lacinia, mauris sed faucibus varius, tellus neque fringilla elit, ut placerat ipsum tellus id ante. Maecenas dapibus dignissim massa, eu pharetra eros hendrerit id. Fusce iaculis dui id est blandit luctus. Morbi eget ornare nunc. In hac habitasse platea dictumst.\n"
//	b := []byte(input)
//
//	actualOutput := DecodeGzip(string(b))
//	assert.Equal(t, expectedOutput, actualOutput)
//
//	// Test the DecodeGzip function with a non-gzip input
//	//nonGzipInput := "This is not a gzip string"
//	//actualOutput = DecodeGzip(nonGzipInput)
//	//assert.Equal(t, nonGzipInput, actualOutput)
//}
