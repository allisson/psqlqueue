// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/allisson/psqlqueue/domain"
	mock "github.com/stretchr/testify/mock"
)

// TopicService is an autogenerated mock type for the TopicService type
type TopicService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, topic
func (_m *TopicService) Create(ctx context.Context, topic *domain.Topic) error {
	ret := _m.Called(ctx, topic)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Topic) error); ok {
		r0 = rf(ctx, topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateMessage provides a mock function with given fields: ctx, topicID, message
func (_m *TopicService) CreateMessage(ctx context.Context, topicID string, message *domain.Message) error {
	ret := _m.Called(ctx, topicID, message)

	if len(ret) == 0 {
		panic("no return value specified for CreateMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *domain.Message) error); ok {
		r0 = rf(ctx, topicID, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *TopicService) Delete(ctx context.Context, id string) error {
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
func (_m *TopicService) Get(ctx context.Context, id string) (*domain.Topic, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *domain.Topic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.Topic, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Topic); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Topic)
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
func (_m *TopicService) List(ctx context.Context, offset uint, limit uint) ([]*domain.Topic, error) {
	ret := _m.Called(ctx, offset, limit)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*domain.Topic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) ([]*domain.Topic, error)); ok {
		return rf(ctx, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) []*domain.Topic); ok {
		r0 = rf(ctx, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Topic)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, uint) error); ok {
		r1 = rf(ctx, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTopicService creates a new instance of TopicService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTopicService(t interface {
	mock.TestingT
	Cleanup(func())
}) *TopicService {
	mock := &TopicService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
