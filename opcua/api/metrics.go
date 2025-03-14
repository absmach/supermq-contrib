// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"time"

	"github.com/absmach/supermq-contrib/opcua"
	"github.com/go-kit/kit/metrics"
)

var _ opcua.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     opcua.Service
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc opcua.Service, counter metrics.Counter, latency metrics.Histogram) opcua.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (mm *metricsMiddleware) CreateClient(ctx context.Context, mgxDevID, opcuaNodeID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_client").Add(1)
		mm.latency.With("method", "create_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateClient(ctx, mgxDevID, opcuaNodeID)
}

func (mm *metricsMiddleware) UpdateClient(ctx context.Context, mgxDevID, opcuaNodeID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "update_client").Add(1)
		mm.latency.With("method", "update_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.UpdateClient(ctx, mgxDevID, opcuaNodeID)
}

func (mm *metricsMiddleware) RemoveClient(ctx context.Context, mgxDevID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "remove_client").Add(1)
		mm.latency.With("method", "remove_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.RemoveClient(ctx, mgxDevID)
}

func (mm *metricsMiddleware) CreateChannel(ctx context.Context, mgxChanID, opcuaServerURI string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_channel").Add(1)
		mm.latency.With("method", "create_channel").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateChannel(ctx, mgxChanID, opcuaServerURI)
}

func (mm *metricsMiddleware) UpdateChannel(ctx context.Context, mgxChanID, opcuaServerURI string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "update_channel").Add(1)
		mm.latency.With("method", "update_channel").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.UpdateChannel(ctx, mgxChanID, opcuaServerURI)
}

func (mm *metricsMiddleware) RemoveChannel(ctx context.Context, mgxChanID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "remove_channel").Add(1)
		mm.latency.With("method", "remove_channel").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.RemoveChannel(ctx, mgxChanID)
}

func (mm *metricsMiddleware) ConnectClient(ctx context.Context, mgxChanID string, mgxClientIDs []string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "connect_client").Add(1)
		mm.latency.With("method", "connect_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.ConnectClient(ctx, mgxChanID, mgxClientIDs)
}

func (mm *metricsMiddleware) DisconnectClient(ctx context.Context, mgxChanID string, mgxClientIDs []string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "disconnect_client").Add(1)
		mm.latency.With("method", "disconnect_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.DisconnectClient(ctx, mgxChanID, mgxClientIDs)
}

func (mm *metricsMiddleware) Browse(ctx context.Context, serverURI, namespace, identifier, identifierType string) ([]opcua.BrowsedNode, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "browse").Add(1)
		mm.latency.With("method", "browse").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Browse(ctx, serverURI, namespace, identifier, identifierType)
}
