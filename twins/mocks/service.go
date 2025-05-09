// Code generated by mockery v2.43.2. DO NOT EDIT.

// Copyright (c) Abstract Machines

package mocks

import (
	context "context"

	messaging "github.com/absmach/supermq/pkg/messaging"
	mock "github.com/stretchr/testify/mock"

	twins "github.com/absmach/supermq-contrib/twins"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// AddTwin provides a mock function with given fields: ctx, token, twin, def
func (_m *Service) AddTwin(ctx context.Context, token string, twin twins.Twin, def twins.Definition) (twins.Twin, error) {
	ret := _m.Called(ctx, token, twin, def)

	if len(ret) == 0 {
		panic("no return value specified for AddTwin")
	}

	var r0 twins.Twin
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, twins.Twin, twins.Definition) (twins.Twin, error)); ok {
		return rf(ctx, token, twin, def)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, twins.Twin, twins.Definition) twins.Twin); ok {
		r0 = rf(ctx, token, twin, def)
	} else {
		r0 = ret.Get(0).(twins.Twin)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, twins.Twin, twins.Definition) error); ok {
		r1 = rf(ctx, token, twin, def)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListStates provides a mock function with given fields: ctx, token, offset, limit, twinID
func (_m *Service) ListStates(ctx context.Context, token string, offset uint64, limit uint64, twinID string) (twins.StatesPage, error) {
	ret := _m.Called(ctx, token, offset, limit, twinID)

	if len(ret) == 0 {
		panic("no return value specified for ListStates")
	}

	var r0 twins.StatesPage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64, uint64, string) (twins.StatesPage, error)); ok {
		return rf(ctx, token, offset, limit, twinID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64, uint64, string) twins.StatesPage); ok {
		r0 = rf(ctx, token, offset, limit, twinID)
	} else {
		r0 = ret.Get(0).(twins.StatesPage)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, uint64, uint64, string) error); ok {
		r1 = rf(ctx, token, offset, limit, twinID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTwins provides a mock function with given fields: ctx, token, offset, limit, name, metadata
func (_m *Service) ListTwins(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata twins.Metadata) (twins.Page, error) {
	ret := _m.Called(ctx, token, offset, limit, name, metadata)

	if len(ret) == 0 {
		panic("no return value specified for ListTwins")
	}

	var r0 twins.Page
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64, uint64, string, twins.Metadata) (twins.Page, error)); ok {
		return rf(ctx, token, offset, limit, name, metadata)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64, uint64, string, twins.Metadata) twins.Page); ok {
		r0 = rf(ctx, token, offset, limit, name, metadata)
	} else {
		r0 = ret.Get(0).(twins.Page)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, uint64, uint64, string, twins.Metadata) error); ok {
		r1 = rf(ctx, token, offset, limit, name, metadata)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveTwin provides a mock function with given fields: ctx, token, twinID
func (_m *Service) RemoveTwin(ctx context.Context, token string, twinID string) error {
	ret := _m.Called(ctx, token, twinID)

	if len(ret) == 0 {
		panic("no return value specified for RemoveTwin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, token, twinID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveStates provides a mock function with given fields: ctx, msg
func (_m *Service) SaveStates(ctx context.Context, msg *messaging.Message) error {
	ret := _m.Called(ctx, msg)

	if len(ret) == 0 {
		panic("no return value specified for SaveStates")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *messaging.Message) error); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTwin provides a mock function with given fields: ctx, token, twin, def
func (_m *Service) UpdateTwin(ctx context.Context, token string, twin twins.Twin, def twins.Definition) error {
	ret := _m.Called(ctx, token, twin, def)

	if len(ret) == 0 {
		panic("no return value specified for UpdateTwin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, twins.Twin, twins.Definition) error); ok {
		r0 = rf(ctx, token, twin, def)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ViewTwin provides a mock function with given fields: ctx, token, twinID
func (_m *Service) ViewTwin(ctx context.Context, token string, twinID string) (twins.Twin, error) {
	ret := _m.Called(ctx, token, twinID)

	if len(ret) == 0 {
		panic("no return value specified for ViewTwin")
	}

	var r0 twins.Twin
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (twins.Twin, error)); ok {
		return rf(ctx, token, twinID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) twins.Twin); ok {
		r0 = rf(ctx, token, twinID)
	} else {
		r0 = ret.Get(0).(twins.Twin)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, token, twinID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
