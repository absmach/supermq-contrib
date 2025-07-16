// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

// Package main contains twins main function to start the twins service.
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
	mongoclient "github.com/absmach/supermq-contrib/pkg/clients/mongo"
	redisclient "github.com/absmach/supermq-contrib/pkg/clients/redis"
	"github.com/absmach/supermq-contrib/twins"
	"github.com/absmach/supermq-contrib/twins/api"
	twapi "github.com/absmach/supermq-contrib/twins/api/http"
	"github.com/absmach/supermq-contrib/twins/events"
	twmongodb "github.com/absmach/supermq-contrib/twins/mongodb"
	"github.com/absmach/supermq-contrib/twins/tracing"
	smqlog "github.com/absmach/supermq/logger"
	"github.com/absmach/supermq/pkg/authn"
	authsvcAuthn "github.com/absmach/supermq/pkg/authn/authsvc"
	"github.com/absmach/supermq/pkg/grpcclient"
	jaegerclient "github.com/absmach/supermq/pkg/jaeger"
	"github.com/absmach/supermq/pkg/messaging"
	"github.com/absmach/supermq/pkg/messaging/brokers"
	brokerstracing "github.com/absmach/supermq/pkg/messaging/brokers/tracing"
	"github.com/absmach/supermq/pkg/prometheus"
	"github.com/absmach/supermq/pkg/server"
	httpserver "github.com/absmach/supermq/pkg/server/http"
	"github.com/absmach/supermq/pkg/uuid"
	"github.com/caarlos0/env/v10"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

const (
	svcName        = "twins"
	envPrefixDB    = "SMQ_TWINS_DB_"
	envPrefixHTTP  = "SMQ_TWINS_HTTP_"
	envPrefixAuth  = "SMQ_AUTH_GRPC_"
	defSvcHTTPPort = "9018"
)

type config struct {
	LogLevel        string  `env:"SMQ_TWINS_LOG_LEVEL"          envDefault:"info"`
	StandaloneID    string  `env:"SMQ_TWINS_STANDALONE_ID"      envDefault:""`
	StandaloneToken string  `env:"SMQ_TWINS_STANDALONE_TOKEN"   envDefault:""`
	ChannelID       string  `env:"SMQ_TWINS_CHANNEL_ID"         envDefault:""`
	BrokerURL       string  `env:"SMQ_MESSAGE_BROKER_URL"       envDefault:"nats://localhost:4222"`
	JaegerURL       url.URL `env:"SMQ_JAEGER_URL"               envDefault:"http://jaeger:14268/api/traces"`
	SendTelemetry   bool    `env:"SMQ_SEND_TELEMETRY"           envDefault:"true"`
	InstanceID      string  `env:"SMQ_TWINS_INSTANCE_ID"        envDefault:""`
	ESURL           string  `env:"SMQ_ES_URL"                   envDefault:"nats://localhost:4222"`
	CacheURL        string  `env:"SMQ_TWINS_CACHE_URL"          envDefault:"redis://localhost:6379/0"`
	TraceRatio      float64 `env:"SMQ_JAEGER_TRACE_RATIO"       envDefault:"1.0"`
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

	cacheClient, err := redisclient.Connect(cfg.CacheURL)
	if err != nil {
		logger.Error(err.Error())
		exitCode = 1
		return
	}
	defer cacheClient.Close()

	db, err := mongoclient.Setup(envPrefixDB)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to setup postgres database : %s", err))
		exitCode = 1
		return
	}

	tp, err := jaegerclient.NewProvider(ctx, svcName, cfg.JaegerURL, cfg.InstanceID, cfg.TraceRatio)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to init Jaeger: %s", err))
		exitCode = 1
		return
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error shutting down tracer provider: %v", err))
		}
	}()
	tracer := tp.Tracer(svcName)

	grpcCfg := grpcclient.Config{}
	if err := env.ParseWithOptions(&grpcCfg, env.Options{Prefix: envPrefixAuth}); err != nil {
		logger.Error(fmt.Sprintf("failed to load auth gRPC client configuration : %s", err))
		exitCode = 1
		return
	}
	authn, authnClient, err := authsvcAuthn.NewAuthentication(ctx, grpcCfg)
	if err != nil {
		logger.Error(err.Error())
		exitCode = 1
		return
	}
	defer authnClient.Close()
	logger.Info("AuthN  successfully connected to auth gRPC server " + authnClient.Secure())

	pubSub, err := brokers.NewPubSub(ctx, cfg.BrokerURL, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to message broker: %s", err))
		exitCode = 1
		return
	}
	defer pubSub.Close()
	pubSub = brokerstracing.NewPubSub(httpServerConfig, tracer, pubSub)

	svc, err := newService(ctx, svcName, pubSub, cfg, authn, tracer, db, cacheClient, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create %s service: %s", svcName, err))
		exitCode = 1
		return
	}

	hs := httpserver.NewServer(ctx, cancel, svcName, httpServerConfig, twapi.MakeHandler(svc, logger, cfg.InstanceID), logger)

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
		logger.Error(fmt.Sprintf("Twins service terminated: %s", err))
	}
}

func newService(ctx context.Context, id string, ps messaging.PubSub, cfg config, authn authn.Authentication, tracer trace.Tracer, db *mongo.Database, cacheclient *redis.Client, logger *slog.Logger) (twins.Service, error) {
	twinRepo := twmongodb.NewTwinRepository(db)
	twinRepo = tracing.TwinRepositoryMiddleware(tracer, twinRepo)

	stateRepo := twmongodb.NewStateRepository(db)
	stateRepo = tracing.StateRepositoryMiddleware(tracer, stateRepo)

	idProvider := uuid.New()
	twinCache := events.NewTwinCache(cacheclient)
	twinCache = tracing.TwinCacheMiddleware(tracer, twinCache)

	svc := twins.New(ps, authn, twinRepo, twinCache, stateRepo, idProvider, cfg.ChannelID, logger)

	var err error
	svc, err = events.NewEventStoreMiddleware(ctx, svc, cfg.ESURL)
	if err != nil {
		return nil, err
	}

	svc = api.LoggingMiddleware(svc, logger)
	counter, latency := prometheus.MakeMetrics(svcName, "api")
	svc = api.MetricsMiddleware(svc, counter, latency)

	subCfg := messaging.SubscriberConfig{
		ID:      id,
		Topic:   brokers.SubjectAllMessages,
		Handler: handle(ctx, logger, cfg.ChannelID, svc),
	}
	if err = ps.Subscribe(ctx, subCfg); err != nil {
		logger.Error(err.Error())
	}

	return svc, nil
}

func handle(ctx context.Context, logger *slog.Logger, chanID string, svc twins.Service) handlerFunc {
	return func(msg *messaging.Message) error {
		if msg.GetChannel() == chanID {
			return nil
		}

		if err := svc.SaveStates(ctx, msg); err != nil {
			logger.Error(fmt.Sprintf("State save failed: %s", err))
			return err
		}

		return nil
	}
}

type handlerFunc func(msg *messaging.Message) error

func (h handlerFunc) Handle(msg *messaging.Message) error {
	return h(msg)
}

func (h handlerFunc) Cancel() error {
	return nil
}
