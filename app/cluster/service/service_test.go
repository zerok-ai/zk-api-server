package service

import (
	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"main/app/cluster/models"
	"main/app/tablemux"
	"main/app/tablemux/mocks"
	"main/app/utils"
	"main/app/utils/zkerrors"
	"px.dev/pxapi"
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

func (s *ServiceTestSuite) TestClusterService_UpdateCluster_ClusterIdEmpty_Created_Success() {
	code, zkErr := s.service.UpdateCluster(nil, models.ClusterDetails{Id: nil})
	assert.Nil(s.T(), zkErr)
	assert.Equal(s.T(), iris.StatusCreated, code)
}

func (s *ServiceTestSuite) TestClusterService_UpdateCluster_ClusterIdNotEmpty_Updated_Success() {
	code, zkErr := s.service.UpdateCluster(nil, models.ClusterDetails{Id: utils.StringToPtr("1")})
	assert.Nil(s.T(), zkErr)
	assert.Equal(s.T(), iris.StatusOK, code)
}

func (s *ServiceTestSuite) TestClusterService_DeleteCluster_ClusterIdNotEmpty_Deleted_Success() {
	code, zkErr := s.service.DeleteCluster(nil, "1")
	assert.Nil(s.T(), zkErr)
	assert.Equal(s.T(), iris.StatusOK, code)
}

func (s *ServiceTestSuite) TestClusterService_List_ResourceDetails_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	clusterIdx, action, st, apiKey := "1", "list", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	pxResp := s.service.GetResourceDetails(nil, clusterIdx, action, st, apiKey)
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), pxResp.Error, &zkErr)
}

func (s *ServiceTestSuite) TestClusterService_Map_ResourceDetails_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	clusterIdx, action, st, apiKey := "1", "map", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetServiceMapMethodSignature(st), DataFrameName: "my_first_map"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	pxResp := s.service.GetResourceDetails(nil, clusterIdx, action, st, apiKey)
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), pxResp.Error, &zkErr)
}

func (s *ServiceTestSuite) TestClusterService_INVALID_ResourceDetails_InternalServerError_Fail() {
	action := "INVALID"
	zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST, "unsupported action: "+action)
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"

	pxResp := s.service.GetResourceDetails(nil, clusterIdx, action, st, apiKey)
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), pxResp.Error, &zkErr)
}

func (s *ServiceTestSuite) TestClusterService_NamespaceList_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
	pxResp := s.service.GetNamespaceList(ctx, "1", "-9MIN", "apiKey")
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), &zkErr, pxResp.Error)
}

func (s *ServiceTestSuite) TestClusterService_NamespaceList_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetNamespaceMethodSignature(st), DataFrameName: "my_first_ns"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	pxResp := s.service.GetNamespaceList(nil, clusterIdx, st, apiKey)
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), pxResp.Error, &zkErr)
}

func (s *ServiceTestSuite) TestClusterService_PodList_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
	pxResp := s.service.GetPodList(ctx, "1", "name", "ns", "-9MIN", "apiKey")
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), &zkErr, pxResp.Error)
}

func (s *ServiceTestSuite) TestClusterService_PodList_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	clusterIdx, name, ns, st, apiKey := "1", "service-name", "namespace", "-9m", "apiKey"
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsMethodSignature(st, ns, ns+"/"+name), DataFrameName: "my_first_graph"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	pxResp := s.service.GetPodList(nil, clusterIdx, name, ns, st, apiKey)
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), pxResp.Error, &zkErr)
}

func (s *ServiceTestSuite) TestClusterService_GetPxlData_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
	pxResp := s.service.GetPxlData(ctx, "1", "-9MIN", "apiKey")
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), &zkErr, pxResp.Error)
}

func (s *ServiceTestSuite) TestClusterService_GetPxlData_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetPXDataSignature(100, st, "{}"), DataFrameName: "my_first_list"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	pxResp := s.service.GetPxlData(nil, clusterIdx, st, apiKey)
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), pxResp.Error, &zkErr)
}

func (s *ServiceTestSuite) TestClusterService_ResourceDetails_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
	pxResp := s.service.GetResourceDetails(ctx, "1", "map", "-9MIN", "apiKey")
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), &zkErr, pxResp.Error)
}

func (s *ServiceTestSuite) TestClusterService_GetServiceDetails_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	clusterIdx, name, ns, st, apiKey := "1", "service-name", "namespace", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetServiceDetailsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	pxResp := s.service.GetServiceDetails(nil, clusterIdx, name, ns, st, apiKey)
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), pxResp.Error, &zkErr)
}

func (s *ServiceTestSuite) TestClusterService_GetServiceDetails_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
	pxResp := s.service.GetServiceDetails(ctx, "1", "service-name", "namespace", "-9MIN", "apiKey")
	assert.Nil(s.T(), pxResp.Result)
	assert.Nil(s.T(), pxResp.ResultsStats)
	assert.Equal(s.T(), &zkErr, pxResp.Error)
}

func TestGetResp(t *testing.T) {
	scriptResults := &pxapi.ScriptResults{}

	resultData := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	expectedOutput1 := map[string]interface{}{
		"results": resultData,
		"stats":   scriptResults.Stats(),
		"status":  200,
	}
	expectedOutput2 := map[string]interface{}{
		"results": nil,
		"stats":   nil,
		"status":  500,
	}

	clusterSvc := &clusterService{}

	// Test case 1: result is non-nil
	actualOutput1 := clusterSvc.getResp(scriptResults, resultData)
	assert.Equal(t, expectedOutput1, actualOutput1)

	// Test case 2: result is nil
	actualOutput2 := clusterSvc.getResp(scriptResults, nil)
	assert.Equal(t, expectedOutput2, actualOutput2)
}
