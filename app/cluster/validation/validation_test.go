package validation

import (
	"github.com/stretchr/testify/assert"
	"main/app/utils"
	"main/app/utils/zkerrors"
	"testing"
)

func TestValidatePxlTime(t *testing.T) {
	// valid time format
	assert.True(t, ValidatePxlTime("1s"))
	assert.True(t, ValidatePxlTime("10m"))

	// invalid time format
	assert.False(t, ValidatePxlTime("s"))            // missing value
	assert.False(t, ValidatePxlTime("1.5s"))         // decimal value
	assert.False(t, ValidatePxlTime("1second"))      // incorrect unit
	assert.False(t, ValidatePxlTime("1second100ms")) // extra values
	assert.True(t, ValidatePxlTime("-1s"))           // negative value
	assert.False(t, ValidatePxlTime("1"))            // missing unit
}

func TestValidateGraphDetailsApi(t *testing.T) {
	// valid input
	zkErr := ValidateGraphDetailsApi("serviceName", "ns", "st", "apiKey")
	assert.Nil(t, zkErr)

	// invalid input
	zkErr = ValidateGraphDetailsApi("", "ns", "st", "apiKey") // empty service name
	assert.NotNil(t, zkErr)
	assert.Equal(t, *zkErr, zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_SERVICE_NAME_EMPTY, nil))

	zkErr = ValidateGraphDetailsApi("serviceName", "", "st", "apiKey") // empty namespace
	assert.NotNil(t, zkErr)
	assert.Equal(t, *zkErr, zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_NAMESPACE_EMPTY, nil))

	zkErr = ValidateGraphDetailsApi("serviceName", "ns", "", "apiKey") // empty time
	assert.NotNil(t, zkErr)
	assert.Equal(t, *zkErr, zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_EMPTY, nil))

	zkErr = ValidateGraphDetailsApi("serviceName", "ns", "st", "") // empty api key
	assert.NotNil(t, zkErr)
	assert.Equal(t, *zkErr, zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_ZK_API_KEY_EMPTY, nil))
}

func TestValidatePodDetailsApi(t *testing.T) {
	testCases := []struct {
		podName, ns, st, apiKey string
		expectedErr             *zkerrors.ZkError
	}{
		// Test case 1: All parameters are valid
		{"pod-1", "namespace-1", "2023-04-25T09:00:00Z", "api-key-1", nil},

		// Test case 2: podName is empty
		{"", "namespace-1", "2023-04-25T09:00:00Z", "api-key-1",
			utils.ToPtr(zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_SERVICE_POD_EMPTY, nil))},

		// Test case 3: ns is empty
		{"pod-1", "", "2023-04-25T09:00:00Z", "api-key-1",
			utils.ToPtr(zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_NAMESPACE_EMPTY, nil))},

		// Test case 4: st is empty
		{"pod-1", "namespace-1", "", "api-key-1",
			utils.ToPtr(zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_EMPTY, nil))},

		// Test case 5: apiKey is empty
		{"pod-1", "namespace-1", "2023-04-25T09:00:00Z", "",
			utils.ToPtr(zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_ZK_API_KEY_EMPTY, nil))},
	}

	for _, tc := range testCases {
		err := ValidatePodDetailsApi(nil, tc.podName, tc.ns, tc.st, tc.apiKey)
		if tc.expectedErr == nil {
			assert.Nil(t, tc.expectedErr)
		} else {
			assert.Equal(t, tc.expectedErr, err)
		}
	}
}

func TestValidateGetResourceDetailsApi(t *testing.T) {
	testCases := []struct {
		st, apiKey  string
		expectedErr *zkerrors.ZkError
	}{
		// Test case 1: All parameters are valid
		{"2023-04-25T09:00:00Z", "api-key-1", nil},

		// Test case 2: st is empty
		{st: "", apiKey: "api-key-1", expectedErr: utils.ToPtr(zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_EMPTY, nil))},

		// Test case 3: apiKey is empty
		{st: "2023-04-25T09:00:00Z", expectedErr: utils.ToPtr(zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_ZK_API_KEY_EMPTY, nil))},
	}

	for _, tc := range testCases {
		err := ValidateGetResourceDetailsApi(tc.st, tc.apiKey)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestValidateGetPxlData(t *testing.T) {
	testCases := []struct {
		s, apiKey   string
		expectedErr *zkerrors.ZkError
	}{
		// Test case 1: All parameters are valid
		{"cluster-1", "api-key-1", nil},

		// Test case 2: s is empty
		{"", "api-key-1",
			utils.ToPtr(zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_CLUSTER_ID_EMPTY, nil))},

		// Test case 3: apiKey is empty
		{"cluster-1", "",
			utils.ToPtr(zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_ZK_API_KEY_EMPTY, nil))},
	}

	for _, tc := range testCases {
		err := ValidateGetPxlData(tc.s, tc.apiKey)
		assert.Equal(t, tc.expectedErr, err)
	}
}