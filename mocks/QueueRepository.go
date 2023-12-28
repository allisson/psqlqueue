// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/allisson/psqlqueue/domain"
	mock "github.com/stretchr/testify/mock"
)

// QueueRepository is an autogenerated mock type for the QueueRepository type
type QueueRepository struct {
	mock.Mock
}

// Cleanup provides a mock function with given fields: ctx, id
func (_m *QueueRepository) Cleanup(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Cleanup")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: ctx, queue
func (_m *QueueRepository) Create(ctx context.Context, queue *domain.Queue) error {
	ret := _m.Called(ctx, queue)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Queue) error); ok {
		r0 = rf(ctx, queue)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *QueueRepository) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, id
func (_m *QueueRepository) Get(ctx context.Context, id string) (*domain.Queue, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *domain.Queue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.Queue, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Queue); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Queue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, offset, limit
func (_m *QueueRepository) List(ctx context.Context, offset int, limit int) ([]*domain.Queue, error) {
	ret := _m.Called(ctx, offset, limit)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*domain.Queue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) ([]*domain.Queue, error)); ok {
		return rf(ctx, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int) []*domain.Queue); ok {
		r0 = rf(ctx, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Queue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Purge provides a mock function with given fields: ctx, id
func (_m *QueueRepository) Purge(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Purge")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stats provides a mock function with given fields: ctx, id
func (_m *QueueRepository) Stats(ctx context.Context, id string) (*domain.QueueStats, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Stats")
	}

	var r0 *domain.QueueStats
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.QueueStats, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.QueueStats); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.QueueStats)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, queue
func (_m *QueueRepository) Update(ctx context.Context, queue *domain.Queue) error {
	ret := _m.Called(ctx, queue)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Queue) error); ok {
		r0 = rf(ctx, queue)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewQueueRepository creates a new instance of QueueRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQueueRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *QueueRepository {
	mock := &QueueRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}