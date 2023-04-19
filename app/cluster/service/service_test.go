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

func (s *ServiceTestSuite) TestClusterService_GetServiceDetails_InternalServerError_Fail() {
	var zkErr zkerrors.ZkError
	zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
	clusterIdx, name, ns, st, apiKey, domain := "1", "service-name", "namespace", "-9m", "apiKey", "px.zkcloud02.getanton.com:443"
	temp := tablemux.MethodTemplate{MethodSignature: utils.GetServiceDetailsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"}

	s.pixieRepoMock.On("GetPixieData", mock.Anything, mock.Anything, temp, clusterIdx, apiKey, domain).Return(nil, &zkErr)

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
