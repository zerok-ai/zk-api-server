package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zerok-ai/zk-utils-go/rules/model"
	"main/app/ruleengine/repository/mocks"
	"main/app/utils"
	"main/app/utils/zkerrors"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	service  RuleService
	repoMock *mocks.RulesRepo
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, &ServiceTestSuite{})
}

// runs before execution of suite
func (s *ServiceTestSuite) SetupSuite() {
	r := mocks.NewRulesRepo()
	s.repoMock = r
	s.service = NewRuleService(r)
}

func (s *ServiceTestSuite) TestRuleService_GetAllRules_twoRules_Success() {

	rs1 := model.RuleSet{
		Rule: model.Rule{
			ID:       "req_method",
			Field:    "req_method",
			Type:     "string",
			Input:    "string",
			Operator: "equal",
			Value:    "POST",
		},
	}

	fr1 := model.FilterRule{
		Version: 1684740000,
		Workloads: map[string]model.WorkloadRule{
			"ws1": {
				Service:   utils.ToPtr[string]("s1"),
				TraceRole: utils.ToPtr[string]("client"),
				Protocol:  utils.ToPtr[string]("HTTP"),
				ConditionalRule: model.ConditionalRule{
					Condition: utils.ToPtr[string]("AND"),
					RuleSet:   []model.RuleSet{rs1},
				},
			},
		},
		FilterId: "f1",
		Filters: model.Filters{
			Type:        "workload",
			Condition:   "AND",
			WorkloadSet: []string{"ws1"},
		},
	}

	rs2 := model.RuleSet{
		Rule: model.Rule{
			ID:       "req_method",
			Field:    "req_method",
			Type:     "string",
			Input:    "string",
			Operator: "equal",
			Value:    "POST",
		},
	}

	fr2 := model.FilterRule{
		Version: 1684749999,
		Workloads: map[string]model.WorkloadRule{
			"ws2": {
				Service:   utils.ToPtr[string]("s2"),
				TraceRole: utils.ToPtr[string]("server"),
				Protocol:  utils.ToPtr[string]("HTTP"),
				ConditionalRule: model.ConditionalRule{
					Condition: utils.ToPtr[string]("AND"),
					RuleSet:   []model.RuleSet{rs2},
				},
			},
		},
		FilterId: "f2",
		Filters: model.Filters{
			Type:        "workload",
			Condition:   "AND",
			WorkloadSet: []string{"ws2"},
		},
	}

	clusterId, version, deleted, offset, limit := "clusterId1", 0, false, 0, 10

	//filter := repository.RuleQueryFilter{
	//	ClusterId: clusterId,
	//	Deleted:   deleted,
	//	Version:   int64(version),
	//	Limit:     limit,
	//	Offset:    offset,
	//}
	// Mock the repository method to return the input JSON
	s.repoMock.On("GetAllRules", mock.Anything).Return(utils.ToPtr([]model.FilterRule{fr1, fr2}), nil).Once()

	res, err := s.service.GetAllRules(clusterId, int64(version), deleted, offset, limit)
	assert.Equal(s.T(), 2, len(res.Rules))
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
	assert.NotNil(s.T(), res.Rules[0].Workloads["ws1"])
	assert.NotNil(s.T(), res.Rules[0].Workloads["ws2"])
}

func (s *ServiceTestSuite) TestRuleService_GetAllRules_RepoErr_Failure() {
	zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_NOT_FOUND, nil)
	s.repoMock.On("GetAllRules", mock.Anything).Return(nil, &zkError).Once()

	res, err := s.service.GetAllRules("clusterId1", 0, false, 0, 10)
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), zkerrors.ZK_ERROR_NOT_FOUND, err.Error)
}
