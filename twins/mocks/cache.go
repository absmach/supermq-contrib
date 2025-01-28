// Code generated by mockery v2.43.2. DO NOT EDIT.

// Copyright (c) Abstract Machines

package mocks

import (
	context "context"

	twins "github.com/absmach/supermq-contrib/twins"
	mock "github.com/stretchr/testify/mock"
)

// TwinCache is an autogenerated mock type for the TwinCache type
type TwinCache struct {
	mock.Mock
}

// IDs provides a mock function with given fields: ctx, channel, subtopic
func (_m *TwinCache) IDs(ctx context.Context, channel string, subtopic string) ([]string, error) {
	ret := _m.Called(ctx, channel, subtopic)

	if len(ret) == 0 {
		panic("no return value specified for IDs")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]string, error)); ok {
		return rf(ctx, channel, subtopic)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []string); ok {
		r0 = rf(ctx, channel, subtopic)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, channel, subtopic)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: ctx, twinID
func (_m *TwinCache) Remove(ctx context.Context, twinID string) error {
	ret := _m.Called(ctx, twinID)

	if len(ret) == 0 {
		panic("no return value specified for Remove")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, twinID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: ctx, twin
func (_m *TwinCache) Save(ctx context.Context, twin twins.Twin) error {
	ret := _m.Called(ctx, twin)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, twins.Twin) error); ok {
		r0 = rf(ctx, twin)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveIDs provides a mock function with given fields: ctx, channel, subtopic, twinIDs
func (_m *TwinCache) SaveIDs(ctx context.Context, channel string, subtopic string, twinIDs []string) error {
	ret := _m.Called(ctx, channel, subtopic, twinIDs)

	if len(ret) == 0 {
		panic("no return value specified for SaveIDs")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, []string) error); ok {
		r0 = rf(ctx, channel, subtopic, twinIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, twin
func (_m *TwinCache) Update(ctx context.Context, twin twins.Twin) error {
	ret := _m.Called(ctx, twin)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, twins.Twin) error); ok {
		r0 = rf(ctx, twin)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTwinCache creates a new instance of TwinCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTwinCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *TwinCache {
	mock := &TwinCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
