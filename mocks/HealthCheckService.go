// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/allisson/psqlqueue/domain"
	mock "github.com/stretchr/testify/mock"
)

// HealthCheckService is an autogenerated mock type for the HealthCheckService type
type HealthCheckService struct {
	mock.Mock
}

// Check provides a mock function with given fields: ctx
func (_m *HealthCheckService) Check(ctx context.Context) (*domain.HealthCheck, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Check")
	}

	var r0 *domain.HealthCheck
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*domain.HealthCheck, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *domain.HealthCheck); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.HealthCheck)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewHealthCheckService creates a new instance of HealthCheckService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHealthCheckService(t interface {
	mock.TestingT
	Cleanup(func())
}) *HealthCheckService {
	mock := &HealthCheckService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
