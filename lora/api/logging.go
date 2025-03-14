// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"log/slog"
	"time"

	"github.com/absmach/supermq-contrib/lora"
)

var _ lora.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    lora.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc lora.Service, logger *slog.Logger) lora.Service {
	return &loggingMiddleware{
		logger: logger,
		svc:    svc,
	}
}

func (lm loggingMiddleware) CreateClient(ctx context.Context, clientID, loraDevEUI string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("client_id", clientID),
			slog.String("dev_eui", loraDevEUI),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Create client route-map failed", args...)
			return
		}
		lm.logger.Info("Create client route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.CreateClient(ctx, clientID, loraDevEUI)
}

func (lm loggingMiddleware) UpdateClient(ctx context.Context, clientID, loraDevEUI string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("client_id", clientID),
			slog.String("dev_eui", loraDevEUI),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update client route-map failed", args...)
			return
		}
		lm.logger.Info("Update client route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateClient(ctx, clientID, loraDevEUI)
}

func (lm loggingMiddleware) RemoveClient(ctx context.Context, clientID string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("client_id", clientID),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Remove client route-map failed", args...)
			return
		}
		lm.logger.Info("Remove client route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.RemoveClient(ctx, clientID)
}

func (lm loggingMiddleware) CreateChannel(ctx context.Context, chanID, loraApp string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", chanID),
			slog.String("lora_app", loraApp),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Create channel route-map failed", args...)
			return
		}
		lm.logger.Info("Create channel route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.CreateChannel(ctx, chanID, loraApp)
}

func (lm loggingMiddleware) UpdateChannel(ctx context.Context, chanID, loraApp string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", chanID),
			slog.String("lora_app", loraApp),
		}
		if err != nil {
			lm.logger.Warn("Update channel route-map failed", args...)
			return
		}
		lm.logger.Info("Update channel route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateChannel(ctx, chanID, loraApp)
}

func (lm loggingMiddleware) RemoveChannel(ctx context.Context, chanID string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", chanID),
		}
		if err != nil {
			lm.logger.Warn("Remove channel route-map failed", args...)
			return
		}
		lm.logger.Info("Remove channel route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.RemoveChannel(ctx, chanID)
}

func (lm loggingMiddleware) ConnectClient(ctx context.Context, chanID, clientID string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", chanID),
			slog.String("client_id", clientID),
		}
		if err != nil {
			args := append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Connect client to channel failed", args...)
			return
		}
		lm.logger.Info("Connect client to channel completed successfully", args...)
	}(time.Now())

	return lm.svc.ConnectClient(ctx, chanID, clientID)
}

func (lm loggingMiddleware) DisconnectClient(ctx context.Context, chanID, clientID string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", chanID),
			slog.String("client_id", clientID),
		}
		if err != nil {
			args := append(args, slog.String("error", err.Error()))
			lm.logger.Warn("Disconnect client from channel failed", args...)
			return
		}
		lm.logger.Info("Disconnect client from channel completed successfully", args...)
	}(time.Now())

	return lm.svc.DisconnectClient(ctx, chanID, clientID)
}

func (lm loggingMiddleware) Publish(ctx context.Context, msg *lora.Message) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("message",
				slog.String("application_id", msg.ApplicationID),
				slog.String("device_eui", msg.DevEUI),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Publish failed", args...)
			return
		}
		lm.logger.Info("Publish completed successfully", args...)
	}(time.Now())

	return lm.svc.Publish(ctx, msg)
}
