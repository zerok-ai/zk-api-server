package utils

// TODO: move the zk-utils-go

//import (
//	"github.com/stretchr/testify/assert"
//	zkhttp "github.com/zerok-ai/zk-utils-go/http"
//	"main/app/utils/errors"
//	"testing"
//)
//
//func TestZkHttpResponse_Header(t *testing.T) {
//	response := zkhttp.ZkHttpResponse[int]{}
//
//	initialHeaders := map[string]string{
//		"Content-Type": "application/json",
//	}
//	response.Headers = &initialHeaders
//
//	response.Header("Authorization", "Bearer token")
//
//	assert.Equal(t, "Bearer token", (*response.Headers)["Authorization"])
//	assert.Equal(t, "application/json", (*response.Headers)["Content-Type"])
//}
//
//func TestZkHttpResponse_Empty_Header(t *testing.T) {
//	response := zkhttp.ZkHttpResponse[int]{}
//	response.Header("Authorization", "Bearer token")
//	assert.Equal(t, "Bearer token", (*response.Headers)["Authorization"])
//}
//
//func TestZkHttpResponse_IsOk(t *testing.T) {
//	response := zkhttp.ZkHttpResponse[string]{
//		Status: 200,
//	}
//
//	assert.True(t, response.IsOk())
//
//	response = zkhttp.ZkHttpResponse[string]{
//		Status: 404,
//	}
//
//	assert.False(t, response.IsOk())
//}
//
//func TestZkHttpResponseBuilder_Debug(t *testing.T) {
//	builder := zkhttp.ZkHttpResponseBuilder[string]{
//		ZkHttpResponse: zkhttp.ZkHttpResponse[string]{},
//	}
//
//	HTTP_DEBUG = true
//	newBuilder := builder.Debug("key", "value")
//
//	assert.Equal(t, builder, *newBuilder)
//
//	assert.NotNil(t, builder.ZkHttpResponse.Debug)
//	assert.Equal(t, "value", (*builder.ZkHttpResponse.Debug)["key"])
//
//	invalidBuilder := _zkHttpResponseBuilder[string]{
//		ZkHttpResponse: ZkHttpResponse[string]{
//			Debug: &map[string]interface{}{},
//		},
//	}
//
//	invalidNewBuilder := invalidBuilder.Debug("key", nil)
//	assert.Equal(t, invalidBuilder, *invalidNewBuilder)
//
//	assert.NotNil(t, invalidBuilder.ZkHttpResponse.Debug)
//	assert.Empty(t, *invalidBuilder.ZkHttpResponse.Debug)
//	HTTP_DEBUG = false
//}
//
//func TestZkHttpError_Build(t *testing.T) {
//	zkError := ZkHttpError{}
//
//	zkErrorType := errors.ZkErrorType{
//		Type:    "error_type",
//		Message: "error_message",
//	}
//	param := "param_value"
//	stack := "stack_trace"
//
//	HTTP_DEBUG = true
//	zkError.Build(zkErrorType, &param, stack)
//
//	assert.Equal(t, "error_type", zkError.Kind)
//	assert.Equal(t, "error_message", zkError.Message)
//	assert.Equal(t, "param_value", zkError.Param)
//	assert.Equal(t, "stack_trace", zkError.Stack)
//
//	zkError = ZkHttpError{}
//
//	HTTP_DEBUG = false
//
//	zkError.Build(zkErrorType, nil, stack)
//
//	assert.Equal(t, "error_type", zkError.Kind)
//	assert.Equal(t, "error_message", zkError.Message)
//	assert.Empty(t, zkError.Param)
//	assert.Nil(t, zkError.Stack)
//}
//
//func TestZkHttpResponseBuilder_WithStatus(t *testing.T) {
//	builder := ZkHttpResponseBuilder[string]{
//		_ZkHttpResponseBuilder: _zkHttpResponseBuilder[string]{},
//	}
//
//	newBuilder := builder.WithStatus(200)
//
//	assert.Equal(t, builder._ZkHttpResponseBuilder, *newBuilder)
//	assert.Equal(t, 200, builder._ZkHttpResponseBuilder.ZkHttpResponse.Status)
//}
//
//func TestZkHttpResponseBuilder_WithZkErrorType(t *testing.T) {
//	builder := ZkHttpResponseBuilder[string]{
//		_ZkHttpResponseBuilder: _zkHttpResponseBuilder[string]{},
//	}
//
//	zkErrorType := errors.ZkErrorType{
//		Type:    "error_type",
//		Message: "error_message",
//		Status:  500,
//	}
//
//	newBuilder := builder.WithZkErrorType(zkErrorType)
//
//	assert.Equal(t, builder._ZkHttpResponseBuilder, *newBuilder)
//	assert.Equal(t, zkErrorType.Status, builder._ZkHttpResponseBuilder.ZkHttpResponse.Status)
//	assert.NotNil(t, builder._ZkHttpResponseBuilder.ZkHttpResponse.Error)
//	assert.Equal(t, "error_type", builder._ZkHttpResponseBuilder.ZkHttpResponse.Error.Kind)
//	assert.Equal(t, "error_message", builder._ZkHttpResponseBuilder.ZkHttpResponse.Error.Message)
//}
//
//func Test_zkHttpResponseBuilder_withStatus(t *testing.T) {
//	builder := _zkHttpResponseBuilder[string]{
//		ZkHttpResponse: ZkHttpResponse[string]{},
//	}
//
//	newBuilder := builder.withStatus(200)
//
//	assert.Equal(t, builder, *newBuilder)
//	assert.Equal(t, 200, builder.ZkHttpResponse.Status)
//}
//
//func Test_zkHttpResponseBuilder_Message(t *testing.T) {
//	builder := _zkHttpResponseBuilder[string]{
//		ZkHttpResponse: ZkHttpResponse[string]{},
//	}
//
//	message := "This is a test message"
//
//	newBuilder := builder.Message(&message)
//
//	assert.Equal(t, builder, *newBuilder)
//
//	assert.Equal(t, &message, builder.ZkHttpResponse.Message)
//}
//
//func Test_zkHttpResponseBuilder_Data(t *testing.T) {
//	builder := _zkHttpResponseBuilder[string]{
//		ZkHttpResponse: ZkHttpResponse[string]{},
//	}
//
//	data := "test data"
//
//	newBuilder := builder.Data(&data)
//
//	assert.Equal(t, builder, *newBuilder)
//
//	assert.Equal(t, "test data", builder.ZkHttpResponse.Data)
//}
//
//func Test_zkHttpResponseBuilder_Debug(t *testing.T) {
//	builder := _zkHttpResponseBuilder[string]{
//		ZkHttpResponse: ZkHttpResponse[string]{},
//	}
//
//	HTTP_DEBUG = true
//
//	debugKey := "debug_key"
//	debugValue := "debug_value"
//
//	newBuilder := builder.Debug(debugKey, debugValue)
//
//	assert.Equal(t, builder, *newBuilder)
//
//	assert.NotNil(t, builder.ZkHttpResponse.Debug)
//	assert.Equal(t, debugValue, (*builder.ZkHttpResponse.Debug)[debugKey])
//}
//
//func Test_zkHttpResponseBuilder_Header(t *testing.T) {
//	builder := _zkHttpResponseBuilder[string]{
//		ZkHttpResponse: ZkHttpResponse[string]{},
//	}
//
//	headerKey := "header_key"
//	headerValue := "header_value"
//
//	newBuilder := builder.Header(headerKey, headerValue)
//
//	assert.Equal(t, builder, *newBuilder)
//
//	assert.NotNil(t, builder.ZkHttpResponse.Headers)
//	assert.Equal(t, headerValue, (*builder.ZkHttpResponse.Headers)[headerKey])
//}
//
//func Test_zkHttpResponseBuilder_Metadata(t *testing.T) {
//	builder := _zkHttpResponseBuilder[string]{
//		ZkHttpResponse: ZkHttpResponse[string]{},
//	}
//
//	metadataKey := "metadata_key"
//	metadataValue := "metadata_value"
//
//	newBuilder := builder.Metadata(metadataKey, metadataValue)
//
//	assert.Equal(t, builder, *newBuilder)
//
//	assert.NotNil(t, builder.ZkHttpResponse.Metadata)
//	assert.Equal(t, metadataValue, (*builder.ZkHttpResponse.Metadata)[metadataKey])
//}
//
//func Test_zkHttpResponseBuilder_Error(t *testing.T) {
//	builder := _zkHttpResponseBuilder[string]{
//		ZkHttpResponse: ZkHttpResponse[string]{},
//	}
//
//	zkHttpError := ZkHttpError{
//		Kind:    "error_kind",
//		Message: "error_message",
//		Param:   "error_param",
//		Stack:   "error_stack",
//		Info:    &map[string]any{},
//	}
//
//	newBuilder := builder.Error(zkHttpError)
//
//	assert.Equal(t, builder, *newBuilder)
//
//	assert.NotNil(t, builder.ZkHttpResponse.Error)
//	assert.Equal(t, zkHttpError, *builder.ZkHttpResponse.Error)
//}
//
//func Test_zkHttpResponseBuilder_ErrorInfo(t *testing.T) {
//	builder := _zkHttpResponseBuilder[string]{
//		ZkHttpResponse: ZkHttpResponse[string]{},
//	}
//
//	errorInfoKey := "error_info_key"
//	errorInfoValue := "error_info_value"
//
//	newBuilder := builder.ErrorInfo(errorInfoKey, errorInfoValue)
//
//	assert.Equal(t, builder, *newBuilder)
//
//	assert.NotNil(t, builder.ZkHttpResponse.Error)
//	assert.NotNil(t, builder.ZkHttpResponse.Error.Info)
//	assert.Equal(t, errorInfoValue, (*builder.ZkHttpResponse.Error.Info)[errorInfoKey])
//}
//
//func Test_zkHttpResponseBuilder_Build(t *testing.T) {
//	builder := _zkHttpResponseBuilder[string]{
//		ZkHttpResponse: ZkHttpResponse[string]{},
//	}
//
//	result := builder.Build()
//
//	assert.Equal(t, builder.ZkHttpResponse, result)
//}
