// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package events

import (
	"context"
	"errors"

	"github.com/absmach/supermq-contrib/lora"
	"github.com/absmach/supermq/pkg/events"
)

const (
	keyType   = "lora"
	keyDevEUI = "dev_eui"
	keyAppID  = "app_id"

	clientPrefix     = "client."
	clientCreate     = clientPrefix + "create"
	clientUpdate     = clientPrefix + "update"
	clientRemove     = clientPrefix + "remove"
	clientConnect    = clientPrefix + "connect"
	clientDisconnect = clientPrefix + "disconnect"

	channelPrefix = "group."
	channelCreate = channelPrefix + "create"
	channelUpdate = channelPrefix + "update"
	channelRemove = channelPrefix + "remove"
)

var (
	errMetadataType = errors.New("field lora is missing in the metadata")

	errMetadataFormat = errors.New("malformed metadata")

	errMetadataAppID = errors.New("application ID not found in channel metadatada")

	errMetadataDevEUI = errors.New("device EUI not found in client metadatada")
)

type eventHandler struct {
	svc lora.Service
}

// NewEventHandler returns new event store handler.
func NewEventHandler(svc lora.Service) events.EventHandler {
	return &eventHandler{
		svc: svc,
	}
}

func (es *eventHandler) Handle(ctx context.Context, event events.Event) error {
	msg, err := event.Encode()
	if err != nil {
		return err
	}

	switch msg["operation"] {
	case clientCreate, clientUpdate:
		cte, derr := decodeCreateClient(msg)
		if derr != nil {
			err = derr
			break
		}
		err = es.svc.CreateClient(ctx, cte.id, cte.loraDevEUI)
	case channelCreate, channelUpdate:
		cce, derr := decodeCreateChannel(msg)
		if derr != nil {
			err = derr
			break
		}
		err = es.svc.CreateChannel(ctx, cce.id, cce.loraAppID)
	case clientRemove:
		rte := decodeRemoveClient(msg)
		err = es.svc.RemoveClient(ctx, rte.id)
	case channelRemove:
		rce := decodeRemoveChannel(msg)
		err = es.svc.RemoveChannel(ctx, rce.id)
	case clientConnect:
		tce := decodeConnectionClient(msg)

		for _, clientID := range tce.clientIDs {
			err = es.svc.ConnectClient(ctx, tce.chanID, clientID)
			if err != nil {
				return err
			}
		}
	case clientDisconnect:
		tde := decodeConnectionClient(msg)

		for _, clientID := range tde.clientIDs {
			err = es.svc.DisconnectClient(ctx, tde.chanID, clientID)
			if err != nil {
				return err
			}
		}
	}
	if err != nil && err != errMetadataType {
		return err
	}

	return nil
}

func decodeCreateClient(event map[string]interface{}) (createClientEvent, error) {
	metadata := events.Read(event, "metadata", map[string]interface{}{})

	cte := createClientEvent{
		id: events.Read(event, "id", ""),
	}

	m, ok := metadata[keyType]
	if !ok {
		return createClientEvent{}, errMetadataType
	}

	lm, ok := m.(map[string]interface{})
	if !ok {
		return createClientEvent{}, errMetadataFormat
	}

	val, ok := lm[keyDevEUI].(string)
	if !ok {
		return createClientEvent{}, errMetadataDevEUI
	}

	cte.loraDevEUI = val
	return cte, nil
}

func decodeRemoveClient(event map[string]interface{}) removeClientEvent {
	return removeClientEvent{
		id: events.Read(event, "id", ""),
	}
}

func decodeCreateChannel(event map[string]interface{}) (createChannelEvent, error) {
	metadata := events.Read(event, "metadata", map[string]interface{}{})

	cce := createChannelEvent{
		id: events.Read(event, "id", ""),
	}

	m, ok := metadata[keyType]
	if !ok {
		return createChannelEvent{}, errMetadataType
	}

	lm, ok := m.(map[string]interface{})
	if !ok {
		return createChannelEvent{}, errMetadataFormat
	}

	val, ok := lm[keyAppID].(string)
	if !ok {
		return createChannelEvent{}, errMetadataAppID
	}

	cce.loraAppID = val
	return cce, nil
}

func decodeConnectionClient(event map[string]interface{}) connectionClientEvent {
	return connectionClientEvent{
		chanID:    events.Read(event, "group_id", ""),
		clientIDs: events.ReadStringSlice(event, "member_ids"),
	}
}

func decodeRemoveChannel(event map[string]interface{}) removeChannelEvent {
	return removeChannelEvent{
		id: events.Read(event, "id", ""),
	}
}
