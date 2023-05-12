// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// RulesRepo is an autogenerated mock type for the RulesRepo type
type RulesRepo struct {
	mock.Mock
}

// GetAllRules provides a mock function with given fields:
func (_m *RulesRepo) GetAllRules() ([]map[string]interface{}, error) {
	ret := _m.Called()

	var r0 []map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]map[string]interface{}, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRulesRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewRulesRepo creates a new instance of RulesRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRulesRepo() *RulesRepo {
	mock := &RulesRepo{}


	return mock
}
