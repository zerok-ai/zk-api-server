package service

import (
	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"main/app/tablemux"
	"main/app/tablemux/mocks"
	"main/app/utils"
	"main/app/utils/errors"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	service       ClusterService
	pixieRepoMock *mocks.PixieRepository
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, &ServiceTestSuite{})
}

// runs before execution of suite
func (s *ServiceTestSuite) SetupSuite() {
	p := mocks.NewPixieRepository()
	s.pixieRepoMock = p
	s.service = NewClusterService(p)
}

func (s *ServiceTestSuite) TestClusterService_List_ResourceDetails_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetServiceDetailsList(nil, clusterIdx, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_List_GetPodDetailsTimeSeries_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, mock.Anything, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetPodDetailsTimeSeries(nil, clusterIdx, "pod_name", "namespace", st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

//func (s *ServiceTestSuite) TestGetServiceDetailsList() {
//	mockRepo := s.pixieRepoMock
//
//	// Create a new clusterService instance with the mock repository
//	cs := &clusterService{pixie: mockRepo}
//
//	// Set up a mock response for the GetPixieData method
//	mockResults := &pxapi.ScriptResults{}
//	mockError := (*errors.ZkError)(nil)
//	mockRepo.On("GetPixieData", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResults, mockError)
//
//	// Create a new Iris context
//	app := iris.New()
//	recorder := httptest.NewRecorder()
//	ctx := app.ContextPool.Acquire(recorder, httptest.NewRequest("GET", "/", nil))
//
//	// Call the GetServiceDetailsList method
//	id := "cluster_id"
//	st := "2022-01-01T00:00:00Z"
//	apiKey := "api_key"
//	resp, err := cs.GetServiceDetailsList(ctx, id, st, apiKey)
//
//	// Verify the mock response
//	assert.Equal(t, mockResults, resp)
//	assert.Equal(t, mockError, err)
//
//	// Verify that the GetPixieData method was called with the expected arguments
//	expectedTx := tablemux.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}
//	mockRepo.AssertCalled(t, "GetPixieData", ctx, mock.AnythingOfType("pxapi.TableMuxer"), expectedTx, id, apiKey, details.Domain)
//}

func (s *ServiceTestSuite) TestClusterService_List_ResourceDetails_IncorrectTime_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
	resp, err := s.service.GetServiceDetailsList(ctx, "1", "-9MIN", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_Map_ResourceDetails_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetServiceMapMethodSignature(st), DataFrameName: "my_first_map"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetServiceDetailsMap(nil, clusterIdx, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestGetPodDetailsTimeSeries_IncorrectTime_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
	resp, err := s.service.GetPodDetailsTimeSeries(ctx, "1", "pod", "ns", "st", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_NamespaceList_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
	resp, err := s.service.GetNamespaceList(ctx, "1", "-9MIN", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_NamespaceList_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetNamespaceMethodSignature(st), DataFrameName: "my_first_ns"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetNamespaceList(nil, clusterIdx, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_PodList_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
	resp, err := s.service.GetPodList(ctx, "1", "name", "ns", "-9MIN", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_PodList_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	clusterIdx, name, ns, st, apiKey := "1", "service-name", "namespace", "-9m", "apiKey"
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsMethodSignature(st, ns, ns+"/"+name), DataFrameName: "my_first_graph"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetPodList(nil, clusterIdx, name, ns, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_GetPxlData_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
	resp, err := s.service.GetPxlData(ctx, "1", "-9MIN", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_GetPxlData_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetPXDataSignature(100, st, "{}"), DataFrameName: "my_first_list"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetPxlData(nil, clusterIdx, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_ResourceDetails_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
	resp, err := s.service.GetServiceDetailsMap(ctx, "1", "-9MIN", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_GetServiceDetails_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
	clusterIdx, name, ns, st, apiKey := "1", "service-name", "namespace", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetServiceDetailsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetServiceDetails(nil, clusterIdx, name, ns, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_GetServiceDetails_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(errors.ZkErrorBadRequestTimeFormat, nil)
	resp, err := s.service.GetServiceDetails(ctx, "1", "service-name", "namespace", "-9MIN", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}
