// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify
// Copyright (c) Abstract Machines

// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"context"

	"github.com/absmach/supermq-contrib/consumers/notifiers"
	"github.com/absmach/supermq/pkg/authn"
	mock "github.com/stretchr/testify/mock"
)

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

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

type Service_Expecter struct {
	mock *mock.Mock
}

func (_m *Service) EXPECT() *Service_Expecter {
	return &Service_Expecter{mock: &_m.Mock}
}

// ConsumeBlocking provides a mock function for the type Service
func (_mock *Service) ConsumeBlocking(ctx context.Context, messages interface{}) error {
	ret := _mock.Called(ctx, messages)

	if len(ret) == 0 {
		panic("no return value specified for ConsumeBlocking")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = returnFunc(ctx, messages)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// Service_ConsumeBlocking_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConsumeBlocking'
type Service_ConsumeBlocking_Call struct {
	*mock.Call
}

// ConsumeBlocking is a helper method to define mock.On call
//   - ctx context.Context
//   - messages interface{}
func (_e *Service_Expecter) ConsumeBlocking(ctx interface{}, messages interface{}) *Service_ConsumeBlocking_Call {
	return &Service_ConsumeBlocking_Call{Call: _e.mock.On("ConsumeBlocking", ctx, messages)}
}

func (_c *Service_ConsumeBlocking_Call) Run(run func(ctx context.Context, messages interface{})) *Service_ConsumeBlocking_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 interface{}
		if args[1] != nil {
			arg1 = args[1].(interface{})
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *Service_ConsumeBlocking_Call) Return(err error) *Service_ConsumeBlocking_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *Service_ConsumeBlocking_Call) RunAndReturn(run func(ctx context.Context, messages interface{}) error) *Service_ConsumeBlocking_Call {
	_c.Call.Return(run)
	return _c
}

// CreateSubscription provides a mock function for the type Service
func (_mock *Service) CreateSubscription(ctx context.Context, session authn.Session, sub notifiers.Subscription) (string, error) {
	ret := _mock.Called(ctx, session, sub)

	if len(ret) == 0 {
		panic("no return value specified for CreateSubscription")
	}

	var r0 string
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, authn.Session, notifiers.Subscription) (string, error)); ok {
		return returnFunc(ctx, session, sub)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, authn.Session, notifiers.Subscription) string); ok {
		r0 = returnFunc(ctx, session, sub)
	} else {
		r0 = ret.Get(0).(string)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, authn.Session, notifiers.Subscription) error); ok {
		r1 = returnFunc(ctx, session, sub)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// Service_CreateSubscription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateSubscription'
type Service_CreateSubscription_Call struct {
	*mock.Call
}

// CreateSubscription is a helper method to define mock.On call
//   - ctx context.Context
//   - session authn.Session
//   - sub notifiers.Subscription
func (_e *Service_Expecter) CreateSubscription(ctx interface{}, session interface{}, sub interface{}) *Service_CreateSubscription_Call {
	return &Service_CreateSubscription_Call{Call: _e.mock.On("CreateSubscription", ctx, session, sub)}
}

func (_c *Service_CreateSubscription_Call) Run(run func(ctx context.Context, session authn.Session, sub notifiers.Subscription)) *Service_CreateSubscription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 authn.Session
		if args[1] != nil {
			arg1 = args[1].(authn.Session)
		}
		var arg2 notifiers.Subscription
		if args[2] != nil {
			arg2 = args[2].(notifiers.Subscription)
		}
		run(
			arg0,
			arg1,
			arg2,
		)
	})
	return _c
}

func (_c *Service_CreateSubscription_Call) Return(s string, err error) *Service_CreateSubscription_Call {
	_c.Call.Return(s, err)
	return _c
}

func (_c *Service_CreateSubscription_Call) RunAndReturn(run func(ctx context.Context, session authn.Session, sub notifiers.Subscription) (string, error)) *Service_CreateSubscription_Call {
	_c.Call.Return(run)
	return _c
}

// ListSubscriptions provides a mock function for the type Service
func (_mock *Service) ListSubscriptions(ctx context.Context, sesssion authn.Session, pm notifiers.PageMetadata) (notifiers.Page, error) {
	ret := _mock.Called(ctx, sesssion, pm)

	if len(ret) == 0 {
		panic("no return value specified for ListSubscriptions")
	}

	var r0 notifiers.Page
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, authn.Session, notifiers.PageMetadata) (notifiers.Page, error)); ok {
		return returnFunc(ctx, sesssion, pm)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, authn.Session, notifiers.PageMetadata) notifiers.Page); ok {
		r0 = returnFunc(ctx, sesssion, pm)
	} else {
		r0 = ret.Get(0).(notifiers.Page)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, authn.Session, notifiers.PageMetadata) error); ok {
		r1 = returnFunc(ctx, sesssion, pm)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// Service_ListSubscriptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListSubscriptions'
type Service_ListSubscriptions_Call struct {
	*mock.Call
}

// ListSubscriptions is a helper method to define mock.On call
//   - ctx context.Context
//   - sesssion authn.Session
//   - pm notifiers.PageMetadata
func (_e *Service_Expecter) ListSubscriptions(ctx interface{}, sesssion interface{}, pm interface{}) *Service_ListSubscriptions_Call {
	return &Service_ListSubscriptions_Call{Call: _e.mock.On("ListSubscriptions", ctx, sesssion, pm)}
}

func (_c *Service_ListSubscriptions_Call) Run(run func(ctx context.Context, sesssion authn.Session, pm notifiers.PageMetadata)) *Service_ListSubscriptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 authn.Session
		if args[1] != nil {
			arg1 = args[1].(authn.Session)
		}
		var arg2 notifiers.PageMetadata
		if args[2] != nil {
			arg2 = args[2].(notifiers.PageMetadata)
		}
		run(
			arg0,
			arg1,
			arg2,
		)
	})
	return _c
}

func (_c *Service_ListSubscriptions_Call) Return(page notifiers.Page, err error) *Service_ListSubscriptions_Call {
	_c.Call.Return(page, err)
	return _c
}

func (_c *Service_ListSubscriptions_Call) RunAndReturn(run func(ctx context.Context, sesssion authn.Session, pm notifiers.PageMetadata) (notifiers.Page, error)) *Service_ListSubscriptions_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveSubscription provides a mock function for the type Service
func (_mock *Service) RemoveSubscription(ctx context.Context, session authn.Session, id string) error {
	ret := _mock.Called(ctx, session, id)

	if len(ret) == 0 {
		panic("no return value specified for RemoveSubscription")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, authn.Session, string) error); ok {
		r0 = returnFunc(ctx, session, id)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// Service_RemoveSubscription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveSubscription'
type Service_RemoveSubscription_Call struct {
	*mock.Call
}

// RemoveSubscription is a helper method to define mock.On call
//   - ctx context.Context
//   - session authn.Session
//   - id string
func (_e *Service_Expecter) RemoveSubscription(ctx interface{}, session interface{}, id interface{}) *Service_RemoveSubscription_Call {
	return &Service_RemoveSubscription_Call{Call: _e.mock.On("RemoveSubscription", ctx, session, id)}
}

func (_c *Service_RemoveSubscription_Call) Run(run func(ctx context.Context, session authn.Session, id string)) *Service_RemoveSubscription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 authn.Session
		if args[1] != nil {
			arg1 = args[1].(authn.Session)
		}
		var arg2 string
		if args[2] != nil {
			arg2 = args[2].(string)
		}
		run(
			arg0,
			arg1,
			arg2,
		)
	})
	return _c
}

func (_c *Service_RemoveSubscription_Call) Return(err error) *Service_RemoveSubscription_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *Service_RemoveSubscription_Call) RunAndReturn(run func(ctx context.Context, session authn.Session, id string) error) *Service_RemoveSubscription_Call {
	_c.Call.Return(run)
	return _c
}

// ViewSubscription provides a mock function for the type Service
func (_mock *Service) ViewSubscription(ctx context.Context, session authn.Session, id string) (notifiers.Subscription, error) {
	ret := _mock.Called(ctx, session, id)

	if len(ret) == 0 {
		panic("no return value specified for ViewSubscription")
	}

	var r0 notifiers.Subscription
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, authn.Session, string) (notifiers.Subscription, error)); ok {
		return returnFunc(ctx, session, id)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, authn.Session, string) notifiers.Subscription); ok {
		r0 = returnFunc(ctx, session, id)
	} else {
		r0 = ret.Get(0).(notifiers.Subscription)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, authn.Session, string) error); ok {
		r1 = returnFunc(ctx, session, id)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// Service_ViewSubscription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ViewSubscription'
type Service_ViewSubscription_Call struct {
	*mock.Call
}

// ViewSubscription is a helper method to define mock.On call
//   - ctx context.Context
//   - session authn.Session
//   - id string
func (_e *Service_Expecter) ViewSubscription(ctx interface{}, session interface{}, id interface{}) *Service_ViewSubscription_Call {
	return &Service_ViewSubscription_Call{Call: _e.mock.On("ViewSubscription", ctx, session, id)}
}

func (_c *Service_ViewSubscription_Call) Run(run func(ctx context.Context, session authn.Session, id string)) *Service_ViewSubscription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 authn.Session
		if args[1] != nil {
			arg1 = args[1].(authn.Session)
		}
		var arg2 string
		if args[2] != nil {
			arg2 = args[2].(string)
		}
		run(
			arg0,
			arg1,
			arg2,
		)
	})
	return _c
}

func (_c *Service_ViewSubscription_Call) Return(subscription notifiers.Subscription, err error) *Service_ViewSubscription_Call {
	_c.Call.Return(subscription, err)
	return _c
}

func (_c *Service_ViewSubscription_Call) RunAndReturn(run func(ctx context.Context, session authn.Session, id string) (notifiers.Subscription, error)) *Service_ViewSubscription_Call {
	_c.Call.Return(run)
	return _c
}
