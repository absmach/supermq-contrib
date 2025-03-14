// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"time"

	"github.com/absmach/supermq-contrib/lora"
	"github.com/go-kit/kit/metrics"
)

var _ lora.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     lora.Service
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc lora.Service, counter metrics.Counter, latency metrics.Histogram) lora.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (mm *metricsMiddleware) CreateClient(ctx context.Context, clientID, loraDevEUI string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_client").Add(1)
		mm.latency.With("method", "create_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateClient(ctx, clientID, loraDevEUI)
}

func (mm *metricsMiddleware) UpdateClient(ctx context.Context, clientID, loraDevEUI string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "update_client").Add(1)
		mm.latency.With("method", "update_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.UpdateClient(ctx, clientID, loraDevEUI)
}

func (mm *metricsMiddleware) RemoveClient(ctx context.Context, clientID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "remove_client").Add(1)
		mm.latency.With("method", "remove_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.RemoveClient(ctx, clientID)
}

func (mm *metricsMiddleware) CreateChannel(ctx context.Context, chanID, loraApp string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_channel").Add(1)
		mm.latency.With("method", "create_channel").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateChannel(ctx, chanID, loraApp)
}

func (mm *metricsMiddleware) UpdateChannel(ctx context.Context, chanID, loraApp string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "update_channel").Add(1)
		mm.latency.With("method", "update_channel").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.UpdateChannel(ctx, chanID, loraApp)
}

func (mm *metricsMiddleware) RemoveChannel(ctx context.Context, chanID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "remove_channel").Add(1)
		mm.latency.With("method", "remove_channel").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.RemoveChannel(ctx, chanID)
}

func (mm *metricsMiddleware) ConnectClient(ctx context.Context, chanID, clientID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "connect_client").Add(1)
		mm.latency.With("method", "connect_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.ConnectClient(ctx, chanID, clientID)
}

func (mm *metricsMiddleware) DisconnectClient(ctx context.Context, chanID, clientID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "disconnect_client").Add(1)
		mm.latency.With("method", "disconnect_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.DisconnectClient(ctx, chanID, clientID)
}

func (mm *metricsMiddleware) Publish(ctx context.Context, msg *lora.Message) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "publish").Add(1)
		mm.latency.With("method", "publish").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Publish(ctx, msg)
}
