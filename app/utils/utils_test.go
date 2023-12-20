package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	zkcommon "github.com/zerok-ai/zk-utils-go/common"
	"math"
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

	var emptySlice []float64
	result = zkcommon.Contains(emptySlice, 5.5)
	assert.False(t, result)
}

func TestIsEmpty(t *testing.T) {
	assert.True(t, zkcommon.IsEmpty(""))

	assert.False(t, zkcommon.IsEmpty("hello"))
}

func TestGetIntegerFromString(t *testing.T) {
	inputStr := "1234"
	expected := 1234
	actual, err := zkcommon.GetIntegerFromString(inputStr)
	assert.Nil(t, err, "Unexpected error for valid input string")
	assert.Equal(t, expected, actual, "Unexpected result for valid input string")

	inputStr = "hello"
	expectedErr := strconv.NumError{
		Func: "Atoi",
		Num:  "hello",
		Err:  errors.New("invalid syntax"),
	}
	actual, err = zkcommon.GetIntegerFromString(inputStr)
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
		result := zkcommon.Round(tc.input, tc.precision)
		if result != tc.expected {
			t.Errorf("Expected Round(%f, %d) to be %f, but got %f", tc.input, tc.precision, tc.expected, result)
		}
	}
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
