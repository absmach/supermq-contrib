// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"

	"github.com/absmach/supermq-contrib/opcua"
	apiutil "github.com/absmach/supermq/api/http/util"
	"github.com/absmach/supermq/pkg/errors"
	"github.com/go-kit/kit/endpoint"
)

func browseEndpoint(svc opcua.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(browseReq)

		if err := req.validate(); err != nil {
			return nil, errors.Wrap(apiutil.ErrValidation, err)
		}

		nodes, err := svc.Browse(ctx, req.ServerURI, req.Namespace, req.Identifier, req.IdentifierType)
		if err != nil {
			return nil, err
		}

		res := browseRes{
			Nodes: nodes,
		}

		return res, nil
	}
}
