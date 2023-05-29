// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import (
	sql "database/sql"

	mock "github.com/stretchr/testify/mock"

	zkerrors "main/app/utils/zkerrors"
)

// DatabaseRepo is an autogenerated mock type for the DatabaseRepo type
type DatabaseRepo[T interface{}] struct {
	mock.Mock
}

// CreateConnection provides a mock function with given fields:
func (_m *DatabaseRepo[T]) CreateConnection() *sql.DB {
	ret := _m.Called()

	var r0 *sql.DB
	if rf, ok := ret.Get(0).(func() *sql.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.DB)
		}
	}

	return r0
}

// Delete provides a mock function with given fields: stmt, param, tx, rollback
func (_m *DatabaseRepo[T]) Delete(stmt string, param []interface{}, tx *sql.Tx, rollback bool) (int, error) {
	ret := _m.Called(stmt, param, tx, rollback)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string, []interface{}, *sql.Tx, bool) (int, error)); ok {
		return rf(stmt, param, tx, rollback)
	}
	if rf, ok := ret.Get(0).(func(string, []interface{}, *sql.Tx, bool) int); ok {
		r0 = rf(stmt, param, tx, rollback)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string, []interface{}, *sql.Tx, bool) error); ok {
		r1 = rf(stmt, param, tx, rollback)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: query, param, args
func (_m *DatabaseRepo[T]) Get(query string, param []interface{}, args ...interface{}) *zkerrors.ZkError {
	var _ca []interface{}
	_ca = append(_ca, query, param)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *zkerrors.ZkError
	if rf, ok := ret.Get(0).(func(string, []interface{}, ...interface{}) *zkerrors.ZkError); ok {
		r0 = rf(query, param, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*zkerrors.ZkError)
		}
	}

	return r0
}

// GetAll provides a mock function with given fields: query, param, rowsProcessor
func (_m *DatabaseRepo[T]) GetAll(query string, param []interface{}, rowsProcessor func(*sql.Rows, error) (*[]T, *[]string, *zkerrors.ZkError)) (*[]T, *[]string, *zkerrors.ZkError) {
	ret := _m.Called(query, param, rowsProcessor)

	var r0 *[]T
	var r1 *[]string
	var r2 *zkerrors.ZkError
	if rf, ok := ret.Get(0).(func(string, []interface{}, func(*sql.Rows, error) (*[]T, *[]string, *zkerrors.ZkError)) (*[]T, *[]string, *zkerrors.ZkError)); ok {
		return rf(query, param, rowsProcessor)
	}
	if rf, ok := ret.Get(0).(func(string, []interface{}, func(*sql.Rows, error) (*[]T, *[]string, *zkerrors.ZkError)) *[]T); ok {
		r0 = rf(query, param, rowsProcessor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]T)
		}
	}

	if rf, ok := ret.Get(1).(func(string, []interface{}, func(*sql.Rows, error) (*[]T, *[]string, *zkerrors.ZkError)) *[]string); ok {
		r1 = rf(query, param, rowsProcessor)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*[]string)
		}
	}

	if rf, ok := ret.Get(2).(func(string, []interface{}, func(*sql.Rows, error) (*[]T, *[]string, *zkerrors.ZkError)) *zkerrors.ZkError); ok {
		r2 = rf(query, param, rowsProcessor)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*zkerrors.ZkError)
		}
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewDatabaseRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewDatabaseRepo creates a new instance of DatabaseRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDatabaseRepo[T interface{}]() *DatabaseRepo[T] {
	mock := &DatabaseRepo[T]{}
	return mock
}
