// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

// Package main contains influxdb-writer main function to start the influxdb-writer service.
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"

	chclient "github.com/absmach/callhome/pkg/client"
	"github.com/absmach/supermq"
	consumertracing "github.com/absmach/supermq-contrib/consumers/tracing"
	"github.com/absmach/supermq-contrib/consumers/writers/api"
	"github.com/absmach/supermq-contrib/consumers/writers/influxdb"
	influxdbclient "github.com/absmach/supermq-contrib/pkg/clients/influxdb"
	"github.com/absmach/supermq/consumers"
	smqlog "github.com/absmach/supermq/logger"
	"github.com/absmach/supermq/pkg/jaeger"
	"github.com/absmach/supermq/pkg/messaging/brokers"
	brokerstracing "github.com/absmach/supermq/pkg/messaging/brokers/tracing"
	"github.com/absmach/supermq/pkg/server"
	httpserver "github.com/absmach/supermq/pkg/server/http"
	"github.com/absmach/supermq/pkg/uuid"
	"github.com/caarlos0/env/v10"
	"golang.org/x/sync/errgroup"
)

const (
	svcName        = "influxdb-writer"
	envPrefixHTTP  = "SMQ_INFLUX_WRITER_HTTP_"
	envPrefixDB    = "SMQ_INFLUXDB_"
	defSvcHTTPPort = "9006"
)

type config struct {
	LogLevel      string  `env:"SMQ_INFLUX_WRITER_LOG_LEVEL"     envDefault:"info"`
	ConfigPath    string  `env:"SMQ_INFLUX_WRITER_CONFIG_PATH"   envDefault:"/config.toml"`
	BrokerURL     string  `env:"SMQ_MESSAGE_BROKER_URL"          envDefault:"nats://localhost:4222"`
	JaegerURL     url.URL `env:"SMQ_JAEGER_URL"                  envDefault:"http://jaeger:14268/api/traces"`
	SendTelemetry bool    `env:"SMQ_SEND_TELEMETRY"              envDefault:"true"`
	InstanceID    string  `env:"SMQ_INFLUX_WRITER_INSTANCE_ID"   envDefault:""`
	TraceRatio    float64 `env:"SMQ_JAEGER_TRACE_RATIO"          envDefault:"1.0"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load %s configuration : %s", svcName, err)
	}

	logger, err := smqlog.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to init logger: %s", err.Error())
	}

	var exitCode int
	defer smqlog.ExitWithError(&exitCode)

	if cfg.InstanceID == "" {
		if cfg.InstanceID, err = uuid.New().ID(); err != nil {
			logger.Error(fmt.Sprintf("failed to generate instanceID: %s", err))
			exitCode = 1
			return
		}
	}

	httpServerConfig := server.Config{Port: defSvcHTTPPort}
	if err := env.ParseWithOptions(&httpServerConfig, env.Options{Prefix: envPrefixHTTP}); err != nil {
		logger.Error(fmt.Sprintf("failed to load %s HTTP server configuration : %s", svcName, err))
		exitCode = 1
		return
	}

	tp, err := jaeger.NewProvider(ctx, svcName, cfg.JaegerURL, cfg.InstanceID, cfg.TraceRatio)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to init Jaeger: %s", err))
		exitCode = 1
		return
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error shutting down tracer provider: %v", err))
		}
	}()
	tracer := tp.Tracer(svcName)

	pubSub, err := brokers.NewPubSub(ctx, cfg.BrokerURL, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to message broker: %s", err))
		exitCode = 1
		return
	}
	defer pubSub.Close()
	pubSub = brokerstracing.NewPubSub(httpServerConfig, tracer, pubSub)

	influxDBConfig := influxdbclient.Config{}
	if err := env.ParseWithOptions(&influxDBConfig, env.Options{Prefix: envPrefixDB}); err != nil {
		logger.Error(fmt.Sprintf("failed to load InfluxDB client configuration from environment variable : %s", err))
		exitCode = 1
		return
	}
	influxDBConfig.DBUrl = fmt.Sprintf("%s://%s:%s", influxDBConfig.Protocol, influxDBConfig.Host, influxDBConfig.Port)

	repocfg := influxdb.RepoConfig{
		Bucket: influxDBConfig.Bucket,
		Org:    influxDBConfig.Org,
	}

	client, err := influxdbclient.Connect(ctx, influxDBConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to InfluxDB : %s", err))
		exitCode = 1
		return
	}
	defer client.Close()

	repo := influxdb.NewAsync(client, repocfg)
	repo = consumertracing.NewAsync(tracer, repo, httpServerConfig)

	// Start consuming and logging errors.
	go func(log *slog.Logger) {
		for err := range repo.Errors() {
			if err != nil {
				log.Error(err.Error())
			}
		}
	}(logger)

	if err := consumers.Start(ctx, svcName, pubSub, repo, cfg.ConfigPath, logger); err != nil {
		logger.Error(fmt.Sprintf("failed to start InfluxDB writer: %s", err))
		exitCode = 1
		return
	}

	hs := httpserver.NewServer(ctx, cancel, svcName, httpServerConfig, api.MakeHandler(svcName, cfg.InstanceID), logger)

	if cfg.SendTelemetry {
		chc := chclient.New(svcName, supermq.Version, logger, cancel)
		go chc.CallHome(ctx)
	}

	g.Go(func() error {
		return hs.Start()
	})

	g.Go(func() error {
		return server.StopSignalHandler(ctx, cancel, logger, svcName, hs)
	})

	if err := g.Wait(); err != nil {
		logger.Error(fmt.Sprintf("InfluxDB reader service terminated: %s", err))
	}
}
