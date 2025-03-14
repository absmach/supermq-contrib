// Code generated by mockery v2.43.2. DO NOT EDIT.

// Copyright (c) Abstract Machines

package mocks

import (
	context "context"

	notifiers "github.com/absmach/supermq-contrib/consumers/notifiers"
	mock "github.com/stretchr/testify/mock"
)

// SubscriptionsRepository is an autogenerated mock type for the SubscriptionsRepository type
type SubscriptionsRepository struct {
	mock.Mock
}

// Remove provides a mock function with given fields: ctx, id
func (_m *SubscriptionsRepository) Remove(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Remove")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Retrieve provides a mock function with given fields: ctx, id
func (_m *SubscriptionsRepository) Retrieve(ctx context.Context, id string) (notifiers.Subscription, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Retrieve")
	}

	var r0 notifiers.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (notifiers.Subscription, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) notifiers.Subscription); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(notifiers.Subscription)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveAll provides a mock function with given fields: ctx, pm
func (_m *SubscriptionsRepository) RetrieveAll(ctx context.Context, pm notifiers.PageMetadata) (notifiers.Page, error) {
	ret := _m.Called(ctx, pm)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveAll")
	}

	var r0 notifiers.Page
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, notifiers.PageMetadata) (notifiers.Page, error)); ok {
		return rf(ctx, pm)
	}
	if rf, ok := ret.Get(0).(func(context.Context, notifiers.PageMetadata) notifiers.Page); ok {
		r0 = rf(ctx, pm)
	} else {
		r0 = ret.Get(0).(notifiers.Page)
	}

	if rf, ok := ret.Get(1).(func(context.Context, notifiers.PageMetadata) error); ok {
		r1 = rf(ctx, pm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, sub
func (_m *SubscriptionsRepository) Save(ctx context.Context, sub notifiers.Subscription) (notifiers.Subscription, error) {
	ret := _m.Called(ctx, sub)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 notifiers.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, notifiers.Subscription) (notifiers.Subscription, error)); ok {
		return rf(ctx, sub)
	}
	if rf, ok := ret.Get(0).(func(context.Context, notifiers.Subscription) notifiers.Subscription); ok {
		r0 = rf(ctx, sub)
	} else {
		r0 = ret.Get(0).(notifiers.Subscription)
	}

	if rf, ok := ret.Get(1).(func(context.Context, notifiers.Subscription) error); ok {
		r1 = rf(ctx, sub)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSubscriptionsRepository creates a new instance of SubscriptionsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSubscriptionsRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *SubscriptionsRepository {
	mock := &SubscriptionsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
