package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	mocks "main/app/ruleengine/repository/mocks"
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

	val := []map[string]interface{}{
		{
			"condition": "AND",
			"zk_request_type": map[string]interface{}{
				"id":       "zk_req_type",
				"field":    "zk_req_type",
				"type":     "string",
				"input":    "string",
				"operator": "equal",
				"value":    "HTTP",
			},
			"rules": []interface{}{
				map[string]interface{}{
					"id":       "zk_req_type",
					"field":    "zk_req_type",
					"type":     "string",
					"input":    "string",
					"operator": "equal/not_equal",
					"value":    "HTTP",
				},
				map[string]interface{}{
					"id":       "req_method",
					"field":    "req_method",
					"type":     "string",
					"input":    "string",
					"operator": "equal",
					"value":    "POST",
				},
				map[string]interface{}{
					"id":       "req_path",
					"field":    "req_path",
					"type":     "string",
					"input":    "string",
					"operator": "ends_with",
					"value":    "/exception",
				},
				map[string]interface{}{
					"id":       "source",
					"field":    "source",
					"type":     "key-value",
					"input":    "key-value",
					"operator": "equal",
					"value":    "{'service_name':'demo/inventory'}",
				},
				map[string]interface{}{
					"id":       "destination",
					"field":    "destination",
					"type":     "key-value",
					"input":    "key-value",
					"operator": "equal",
					"value":    "{'service_name':'demo/inventory'}",
				},
			},
			"valid": true,
		},
		{
			"condition": "OR",
			"zk_request_type": map[string]interface{}{
				"id":       "zk_req_type",
				"field":    "zk_req_type",
				"type":     "string",
				"input":    "string",
				"operator": "equal",
				"value":    "HTTP",
			},
			"rules": []interface{}{
				map[string]interface{}{
					"id":       "zk_req_type",
					"field":    "zk_req_type",
					"type":     "string",
					"input":    "string",
					"operator": "equal/not_equal",
					"value":    "HTTP",
				},
				map[string]interface{}{
					"id":       "req_method",
					"field":    "req_method",
					"type":     "string",
					"input":    "string",
					"operator": "equal",
					"value":    "POST",
				},
				map[string]interface{}{
					"id":       "req_path",
					"field":    "req_path",
					"type":     "string",
					"input":    "string",
					"operator": "ends_with",
					"value":    "/exception",
				},
				map[string]interface{}{
					"id":       "source",
					"field":    "source",
					"type":     "key-value",
					"input":    "key-value",
					"operator": "equal",
					"value":    "{'service_name':'demo/inventory'}",
				},
				map[string]interface{}{
					"id":       "destination",
					"field":    "destination",
					"type":     "key-value",
					"input":    "key-value",
					"operator": "equal",
					"value":    "{'service_name':'demo/inventory'}",
				},
			},
			"valid": true,
		},
	}

	// Mock the repository method to return the input JSON
	s.repoMock.On("GetAllRules").Return(val, nil).Once()

	res, err := s.service.GetAllRules()
	assert.Equal(s.T(), 2, len(res.Rules))
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
	assert.Equal(s.T(), "AND", *(res.Rules[0].Condition))
	assert.Equal(s.T(), "OR", *(res.Rules[1].Condition))
}

func (s *ServiceTestSuite) TestRuleService_GetAllRules_UnmarshallErr_Failure() {

	val := []map[string]interface{}{
		{
			"abc": "AND",
		},
		{
			"xyz": 1,
		},
		{
			"jkl": true,
		},
	}

	// Mock the repository method to return the input JSON
	s.repoMock.On("GetAllRules").Return(val, nil).Once()

	res, err := s.service.GetAllRules()
	assert.Equal(s.T(), zkerrors.ZK_ERROR_INTERNAL_SERVER, err.Error)
	assert.Nil(s.T(), res)
}

func (s *ServiceTestSuite) TestRuleService_GetAllRules_RepoErr_Failure() {
	s.repoMock.On("GetAllRules").Return(nil, errors.New("some err from repo")).Once()

	res, err := s.service.GetAllRules()
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), zkerrors.ZK_ERROR_INTERNAL_SERVER, err.Error)
}
