// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"

	notifiers "github.com/absmach/supermq-contrib/consumers/notifiers"
	api "github.com/absmach/supermq/api/http"
	apiutil "github.com/absmach/supermq/api/http/util"
	"github.com/absmach/supermq/pkg/authn"
	"github.com/absmach/supermq/pkg/errors"
	svcerr "github.com/absmach/supermq/pkg/errors/service"
	"github.com/go-kit/kit/endpoint"
)

func createSubscriptionEndpoint(svc notifiers.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createSubReq)
		if err := req.validate(); err != nil {
			return createSubRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}
		session, ok := ctx.Value(api.SessionKey).(authn.Session)
		if !ok {
			return nil, svcerr.ErrAuthentication
		}
		sub := notifiers.Subscription{
			Contact: req.Contact,
			Topic:   req.Topic,
		}
		newSub, err := svc.CreateSubscription(ctx, session, sub)
		if err != nil {
			return createSubRes{}, err
		}

		ucr := createSubRes{
			Subscription: newSub,
		}

		return ucr, nil
	}
}

func viewSubscriptionEndpint(svc notifiers.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(subReq)
		if err := req.validate(); err != nil {
			return viewSubRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}
		session, ok := ctx.Value(api.SessionKey).(authn.Session)
		if !ok {
			return nil, svcerr.ErrAuthentication
		}
		sub, err := svc.ViewSubscription(ctx, session, req.id)
		if err != nil {
			return viewSubRes{}, err
		}
		res := viewSubRes{
			Subscription: sub,
		}
		return res, nil
	}
}

func listSubscriptionsEndpoint(svc notifiers.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listSubsReq)
		if err := req.validate(); err != nil {
			return listSubsRes{}, errors.Wrap(apiutil.ErrValidation, err)
		}
		session, ok := ctx.Value(api.SessionKey).(authn.Session)
		if !ok {
			return nil, svcerr.ErrAuthentication
		}
		pm := notifiers.PageMetadata{
			Topic:   req.topic,
			Contact: req.contact,
			Offset:  req.offset,
			Limit:   int(req.limit),
		}
		page, err := svc.ListSubscriptions(ctx, session, pm)
		if err != nil {
			return listSubsRes{}, err
		}
		res := listSubsRes{
			Page: page,
		}

		return res, nil
	}
}

func deleteSubscriptionEndpint(svc notifiers.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(subReq)
		if err := req.validate(); err != nil {
			return nil, errors.Wrap(apiutil.ErrValidation, err)
		}

		session, ok := ctx.Value(api.SessionKey).(authn.Session)
		if !ok {
			return nil, svcerr.ErrAuthentication
		}
		if err := svc.RemoveSubscription(ctx, session, req.id); err != nil {
			return nil, err
		}
		return removeSubRes{}, nil
	}
}
