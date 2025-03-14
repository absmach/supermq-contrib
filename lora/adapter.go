// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package lora

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/absmach/supermq/pkg/messaging"
)

const protocol = "lora"

var (
	// ErrMalformedMessage indicates malformed LoRa message.
	ErrMalformedMessage = errors.New("malformed message received")

	// ErrNotFoundDev indicates a non-existent route map for a device EUI.
	ErrNotFoundDev = errors.New("route map not found for this device EUI")

	// ErrNotFoundApp indicates a non-existent route map for an application ID.
	ErrNotFoundApp = errors.New("route map not found for this application ID")

	// ErrNotConnected indicates a non-existent route map for a connection.
	ErrNotConnected = errors.New("route map not found for this connection")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// CreateClient creates clientID:devEUI route-map
	CreateClient(ctx context.Context, clientID, devEUI string) error

	// UpdateClient updates clientID:devEUI route-map
	UpdateClient(ctx context.Context, clientID, devEUI string) error

	// RemoveClient removes clientID:devEUI route-map
	RemoveClient(ctx context.Context, clientID string) error

	// CreateChannel creates channelID:appID route-map
	CreateChannel(ctx context.Context, chanID, appID string) error

	// UpdateChannel updates channelID:appID route-map
	UpdateChannel(ctx context.Context, chanID, appID string) error

	// RemoveChannel removes channelID:appID route-map
	RemoveChannel(ctx context.Context, chanID string) error

	// ConnectClient creates clientID:channelID route-map
	ConnectClient(ctx context.Context, chanID, clientID string) error

	// DisconnectClient removes clientID:channelID route-map
	DisconnectClient(ctx context.Context, chanID, clientID string) error

	// Publish forwards messages from the LoRa MQTT broker to SupeMQ Message Broker
	Publish(ctx context.Context, msg *Message) error
}

var _ Service = (*adapterService)(nil)

type adapterService struct {
	publisher  messaging.Publisher
	clientsRM  RouteMapRepository
	channelsRM RouteMapRepository
	connectRM  RouteMapRepository
}

// New instantiates the LoRa adapter implementation.
func New(publisher messaging.Publisher, clientsRM, channelsRM, connectRM RouteMapRepository) Service {
	return &adapterService{
		publisher:  publisher,
		clientsRM:  clientsRM,
		channelsRM: channelsRM,
		connectRM:  connectRM,
	}
}

// Publish forwards messages from Lora MQTT broker to SupeMQ Message broker.
func (as *adapterService) Publish(ctx context.Context, m *Message) error {
	// Get route map of lora application
	clientID, err := as.clientsRM.Get(ctx, m.DevEUI)
	if err != nil {
		return ErrNotFoundDev
	}

	// Get route map of lora application
	chanID, err := as.channelsRM.Get(ctx, m.ApplicationID)
	if err != nil {
		return ErrNotFoundApp
	}

	c := fmt.Sprintf("%s:%s", chanID, clientID)
	if _, err := as.connectRM.Get(ctx, c); err != nil {
		return ErrNotConnected
	}

	// Use the SenML message decoded on LoRa Server application if
	// field Object isn't empty. Otherwise, decode standard field Data.
	var payload []byte
	switch m.Object {
	case nil:
		payload, err = base64.StdEncoding.DecodeString(m.Data)
		if err != nil {
			return ErrMalformedMessage
		}
	default:
		jo, err := json.Marshal(m.Object)
		if err != nil {
			return err
		}
		payload = jo
	}

	// Publish on SupeMQ Message broker
	msg := messaging.Message{
		Publisher: clientID,
		Protocol:  protocol,
		Channel:   chanID,
		Payload:   payload,
		Created:   time.Now().UnixNano(),
	}

	return as.publisher.Publish(ctx, msg.Channel, &msg)
}

func (as *adapterService) CreateClient(ctx context.Context, clientID, devEUI string) error {
	return as.clientsRM.Save(ctx, clientID, devEUI)
}

func (as *adapterService) UpdateClient(ctx context.Context, clientID, devEUI string) error {
	return as.clientsRM.Save(ctx, clientID, devEUI)
}

func (as *adapterService) RemoveClient(ctx context.Context, clientID string) error {
	return as.clientsRM.Remove(ctx, clientID)
}

func (as *adapterService) CreateChannel(ctx context.Context, chanID, appID string) error {
	return as.channelsRM.Save(ctx, chanID, appID)
}

func (as *adapterService) UpdateChannel(ctx context.Context, chanID, appID string) error {
	return as.channelsRM.Save(ctx, chanID, appID)
}

func (as *adapterService) RemoveChannel(ctx context.Context, chanID string) error {
	return as.channelsRM.Remove(ctx, chanID)
}

func (as *adapterService) ConnectClient(ctx context.Context, chanID, clientID string) error {
	if _, err := as.channelsRM.Get(ctx, chanID); err != nil {
		return ErrNotFoundApp
	}

	if _, err := as.clientsRM.Get(ctx, clientID); err != nil {
		return ErrNotFoundDev
	}

	c := fmt.Sprintf("%s:%s", chanID, clientID)
	return as.connectRM.Save(ctx, c, c)
}

func (as *adapterService) DisconnectClient(ctx context.Context, chanID, clientID string) error {
	if _, err := as.channelsRM.Get(ctx, chanID); err != nil {
		return ErrNotFoundApp
	}

	if _, err := as.clientsRM.Get(ctx, clientID); err != nil {
		return ErrNotFoundDev
	}

	c := fmt.Sprintf("%s:%s", chanID, clientID)
	return as.connectRM.Remove(ctx, c)
}
