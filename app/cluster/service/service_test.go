package service

import (
	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"main/app/tablemux"
	"main/app/tablemux/mocks"
	"main/app/utils"
	"main/app/utils/zkerrors"
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
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetServiceDetailsList(nil, clusterIdx, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_Map_ResourceDetails_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetServiceMapMethodSignature(st), DataFrameName: "my_first_map"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetServiceDetailsMap(nil, clusterIdx, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_NamespaceList_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
	resp, err := s.service.GetNamespaceList(ctx, "1", "-9MIN", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_NamespaceList_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetNamespaceMethodSignature(st), DataFrameName: "my_first_ns"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetNamespaceList(nil, clusterIdx, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_PodList_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
	resp, err := s.service.GetPodList(ctx, "1", "name", "ns", "-9MIN", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_PodList_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	clusterIdx, name, ns, st, apiKey := "1", "service-name", "namespace", "-9m", "apiKey"
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetPodDetailsMethodSignature(st, ns, ns+"/"+name), DataFrameName: "my_first_graph"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetPodList(nil, clusterIdx, name, ns, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_GetPxlData_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
	resp, err := s.service.GetPxlData(ctx, "1", "-9MIN", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_GetPxlData_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	clusterIdx, st, apiKey := "1", "-9m", "apiKey"
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetPXDataSignature(100, st, "{}"), DataFrameName: "my_first_list"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetPxlData(nil, clusterIdx, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_ResourceDetails_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
	resp, err := s.service.GetServiceDetailsMap(ctx, "1", "-9MIN", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_GetServiceDetails_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	clusterIdx, name, ns, st, apiKey := "1", "service-name", "namespace", "-9m", "apiKey"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetServiceDetailsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, mock.Anything).Return(nil, &zkErr)

	resp, err := s.service.GetServiceDetails(nil, clusterIdx, name, ns, st, apiKey)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}

func (s *ServiceTestSuite) TestClusterService_GetServiceDetails_IncorrectTimeFormat_Fail() {
	var ctx iris.Context
	zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_BAD_REQUEST_TIME_FORMAT, nil)
	resp, err := s.service.GetServiceDetails(ctx, "1", "service-name", "namespace", "-9MIN", "apiKey")
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), &zkErr, err)
}
