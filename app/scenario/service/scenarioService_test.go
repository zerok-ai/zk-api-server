package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"main/app/scenario/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
)

type ServiceTestSuite struct {
	suite.Suite
	service  ScenarioService
	repoMock *mocks.ScenarioRepo
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, &ServiceTestSuite{})
}

// runs before execution of suite
func (s *ServiceTestSuite) SetupSuite() {
	r := mocks.NewScenarioRepo()
	s.repoMock = r
	s.service = NewScenarioService(r)
}

func (s *ServiceTestSuite) TestScenarioService_GetAllScenario_oneScenarios_Success() {

	var s1 model.Scenario
	validScenarioJsonString := string(zkCommon.GetBytesFromFile("files/validScenarioJsonString.json"))
	err1 := json.Unmarshal([]byte(validScenarioJsonString), &s1)
	assert.NoError(s.T(), err1)

	clusterId, version, deleted, offset, limit := "clusterId1", 0, false, 0, 10
	s.repoMock.On("GetAllScenario", mock.Anything).Return(zkCommon.ToPtr([]model.Scenario{s1}), zkCommon.ToPtr([]string{"deleted_scenario_id"}), nil).Once()

	res, err := s.service.GetAllScenario(clusterId, int64(version), deleted, offset, limit)
	assert.Equal(s.T(), 1, len(res.Scenarios))
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
	assert.NotNil(s.T(), (*res.Scenarios[0].Workloads)["idA"])
	_, y := (*res.Scenarios[0].Workloads)["any_random_key"]
	assert.False(s.T(), y)
}

func (s *ServiceTestSuite) TestScenarioService_GetAllScenarios_RepoErr_Failure() {
	zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorNotFound, nil)
	s.repoMock.On("GetAllScenario", mock.Anything).Return(nil, nil, sql.ErrNoRows).Once()

	res, err := s.service.GetAllScenario("clusterId1", 0, false, 0, 10)
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), zkError.Error, err.Error)
}

func (s *ServiceTestSuite) TestScenarioService_GetAllScenarios_SomeErr_Failure() {
	zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
	s.repoMock.On("GetAllScenario", mock.Anything).Return(nil, nil, errors.New("some err")).Once()

	res, err := s.service.GetAllScenario("clusterId1", 0, false, 0, 10)
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), zkError.Error, err.Error)
}
