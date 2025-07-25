// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

// Package main contains lora main function to start the lora service.
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"time"

	chclient "github.com/absmach/callhome/pkg/client"
	"github.com/absmach/supermq"
	"github.com/absmach/supermq-contrib/lora"
	"github.com/absmach/supermq-contrib/lora/api"
	loraevents "github.com/absmach/supermq-contrib/lora/events"
	"github.com/absmach/supermq-contrib/lora/mqtt"
	redisclient "github.com/absmach/supermq-contrib/pkg/clients/redis"
	smqlog "github.com/absmach/supermq/logger"
	"github.com/absmach/supermq/pkg/events"
	"github.com/absmach/supermq/pkg/events/store"
	"github.com/absmach/supermq/pkg/jaeger"
	"github.com/absmach/supermq/pkg/messaging"
	"github.com/absmach/supermq/pkg/messaging/brokers"
	brokerstracing "github.com/absmach/supermq/pkg/messaging/brokers/tracing"
	"github.com/absmach/supermq/pkg/prometheus"
	"github.com/absmach/supermq/pkg/server"
	httpserver "github.com/absmach/supermq/pkg/server/http"
	"github.com/absmach/supermq/pkg/uuid"
	"github.com/caarlos0/env/v10"
	mqttpaho "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/errgroup"
)

const (
	svcName        = "lora-adapter"
	envPrefixHTTP  = "SMQ_LORA_ADAPTER_HTTP_"
	defSvcHTTPPort = "9017"

	clientsRMPrefix  = "client"
	channelsRMPrefix = "channel"
	connsRMPrefix    = "connection"
	clientsStream    = "events.supermq.client.*"
	channelsStream   = "events.supermq.channel.*"
)

type config struct {
	LogLevel       string        `env:"SMQ_LORA_ADAPTER_LOG_LEVEL"           envDefault:"info"`
	LoraMsgURL     string        `env:"SMQ_LORA_ADAPTER_MESSAGES_URL"        envDefault:"tcp://localhost:1883"`
	LoraMsgUser    string        `env:"SMQ_LORA_ADAPTER_MESSAGES_USER"       envDefault:""`
	LoraMsgPass    string        `env:"SMQ_LORA_ADAPTER_MESSAGES_PASS"       envDefault:""`
	LoraMsgTopic   string        `env:"SMQ_LORA_ADAPTER_MESSAGES_TOPIC"      envDefault:"application/+/device/+/event/up"`
	LoraMsgTimeout time.Duration `env:"SMQ_LORA_ADAPTER_MESSAGES_TIMEOUT"    envDefault:"30s"`
	ESConsumerName string        `env:"SMQ_LORA_ADAPTER_EVENT_CONSUMER"      envDefault:"lora-adapter"`
	BrokerURL      string        `env:"SMQ_MESSAGE_BROKER_URL"               envDefault:"nats://localhost:4222"`
	JaegerURL      url.URL       `env:"SMQ_JAEGER_URL"                       envDefault:"http://localhost:14268/api/traces"`
	SendTelemetry  bool          `env:"SMQ_SEND_TELEMETRY"                   envDefault:"true"`
	InstanceID     string        `env:"SMQ_LORA_ADAPTER_INSTANCE_ID"         envDefault:""`
	ESURL          string        `env:"SMQ_ES_URL"                           envDefault:"nats://localhost:4222"`
	RouteMapURL    string        `env:"SMQ_LORA_ADAPTER_ROUTE_MAP_URL"       envDefault:"redis://localhost:6379/0"`
	TraceRatio     float64       `env:"SMQ_JAEGER_TRACE_RATIO"               envDefault:"1.0"`
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

	rmConn, err := redisclient.Connect(cfg.RouteMapURL)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to setup route map redis client : %s", err))
		exitCode = 1
		return
	}
	defer rmConn.Close()

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

	pub, err := brokers.NewPublisher(ctx, cfg.BrokerURL)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to message broker: %s", err))
		exitCode = 1
		return
	}
	defer pub.Close()
	pub = brokerstracing.NewPublisher(httpServerConfig, tracer, pub)

	svc := newService(pub, rmConn, clientsRMPrefix, channelsRMPrefix, connsRMPrefix, logger)

	mqttConn, err := connectToMQTTBroker(cfg.LoraMsgURL, cfg.LoraMsgUser, cfg.LoraMsgPass, cfg.LoraMsgTimeout, logger)
	if err != nil {
		logger.Error(err.Error())
		exitCode = 1
		return
	}

	if err = subscribeToLoRaBroker(svc, mqttConn, cfg.LoraMsgTimeout, cfg.LoraMsgTopic, logger); err != nil {
		logger.Error(fmt.Sprintf("failed to subscribe to Lora MQTT broker: %s", err))
		exitCode = 1
		return
	}

	if err = subscribeToClientsES(ctx, svc, cfg, logger); err != nil {
		logger.Error(fmt.Sprintf("failed to subscribe to clients event store: %s", err))
		exitCode = 1
		return
	}

	if err = subscribeToChannelsES(ctx, svc, cfg, logger); err != nil {
		logger.Error(fmt.Sprintf("failed to subscribe to channels event store: %s", err))
		exitCode = 1
		return
	}

	logger.Info("Subscribed to Event Store")

	hs := httpserver.NewServer(ctx, cancel, svcName, httpServerConfig, api.MakeHandler(cfg.InstanceID), logger)

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
		logger.Error(fmt.Sprintf("LoRa adapter terminated: %s", err))
	}
}

func connectToMQTTBroker(burl, user, password string, timeout time.Duration, logger *slog.Logger) (mqttpaho.Client, error) {
	opts := mqttpaho.NewClientOptions()
	opts.AddBroker(burl)
	opts.SetClientID(svcName)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetOnConnectHandler(func(_ mqttpaho.Client) {
		logger.Info("Connected to Lora MQTT broker")
	})
	opts.SetConnectionLostHandler(func(_ mqttpaho.Client, err error) {
		logger.Error(fmt.Sprintf("MQTT connection lost: %s", err))
	})

	client := mqttpaho.NewClient(opts)

	if token := client.Connect(); token.WaitTimeout(timeout) && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to Lora MQTT broker: %s", token.Error())
	}

	return client, nil
}

func subscribeToLoRaBroker(svc lora.Service, mc mqttpaho.Client, timeout time.Duration, topic string, logger *slog.Logger) error {
	mqttBroker := mqtt.NewBroker(svc, mc, timeout, logger)
	logger.Info("Subscribed to Lora MQTT broker")
	if err := mqttBroker.Subscribe(topic); err != nil {
		return fmt.Errorf("failed to subscribe to Lora MQTT broker: %s", err)
	}
	return nil
}

func subscribeToClientsES(ctx context.Context, svc lora.Service, cfg config, logger *slog.Logger) error {
	subscriber, err := store.NewSubscriber(ctx, cfg.ESURL, logger)
	if err != nil {
		return err
	}

	subConfig := events.SubscriberConfig{
		Stream:   clientsStream,
		Consumer: cfg.ESConsumerName,
		Handler:  loraevents.NewEventHandler(svc),
	}
	return subscriber.Subscribe(ctx, subConfig)
}

func subscribeToChannelsES(ctx context.Context, svc lora.Service, cfg config, logger *slog.Logger) error {
	subscriber, err := store.NewSubscriber(ctx, cfg.ESURL, logger)
	if err != nil {
		return err
	}

	subConfig := events.SubscriberConfig{
		Stream:   channelsStream,
		Consumer: cfg.ESConsumerName,
		Handler:  loraevents.NewEventHandler(svc),
	}
	return subscriber.Subscribe(ctx, subConfig)
}

func newRouteMapRepository(client *redis.Client, prefix string, logger *slog.Logger) lora.RouteMapRepository {
	logger.Info(fmt.Sprintf("Connected to %s Redis Route-map", prefix))
	return loraevents.NewRouteMapRepository(client, prefix)
}

func newService(pub messaging.Publisher, rmConn *redis.Client, clientsRMPrefix, channelsRMPrefix, connsRMPrefix string, logger *slog.Logger) lora.Service {
	clientsRM := newRouteMapRepository(rmConn, clientsRMPrefix, logger)
	chansRM := newRouteMapRepository(rmConn, channelsRMPrefix, logger)
	connsRM := newRouteMapRepository(rmConn, connsRMPrefix, logger)

	svc := lora.New(pub, clientsRM, chansRM, connsRM)
	svc = api.LoggingMiddleware(svc, logger)
	counter, latency := prometheus.MakeMetrics("lora_adapter", "api")
	svc = api.MetricsMiddleware(svc, counter, latency)

	return svc
}
