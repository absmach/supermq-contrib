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

	clientPrefix = "client."
	clientCreate = clientPrefix + "create"
	clientUpdate = clientPrefix + "update"
	clientRemove = clientPrefix + "remove"

	channelPrefix     = "channel."
	channelCreate     = channelPrefix + "create"
	channelUpdate     = channelPrefix + "update"
	channelRemove     = channelPrefix + "remove"
	channelConnect    = channelPrefix + "connect"
	channelDisconnect = channelPrefix + "disconnect"
)

var (
	errMetadataType = errors.New("field lora is missing in the metadata")

	errMetadataFormat = errors.New("malformed metadata")

	errMetadataAppID = errors.New("application ID not found in channel metadata")

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
		err = es.svc.CreateChannel(ctx, cce.channelID, cce.domainID, cce.loraAppID)
	case clientRemove:
		rte := decodeRemoveClient(msg)
		err = es.svc.RemoveClient(ctx, rte.id)
	case channelRemove:
		rce := decodeRemoveChannel(msg)
		err = es.svc.RemoveChannel(ctx, rce.channelID, rce.domainID)
	case channelConnect:
		tce := decodeConnection(msg)

		for _, chanID := range tce.chanIDs {
			for _, clientID := range tce.clientIDs {
				err = es.svc.ConnectClient(ctx, chanID, tce.domainID, clientID)
				if err != nil {
					return err
				}
			}
		}
	case channelDisconnect:
		tde := decodeConnection(msg)

		for _, chanID := range tde.chanIDs {
			for _, clientID := range tde.clientIDs {
				err = es.svc.DisconnectClient(ctx, chanID, tde.domainID, clientID)
				if err != nil {
					return err
				}
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
		channelID: events.Read(event, "id", ""),
		domainID:  events.Read(event, "domain", ""),
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

func decodeConnection(event map[string]interface{}) connectionEvent {
	return connectionEvent{
		chanIDs:   events.ReadStringSlice(event, "channel_ids"),
		clientIDs: events.ReadStringSlice(event, "client_ids"),
		domainID:  events.Read(event, "domain", ""),
	}
}

func decodeRemoveChannel(event map[string]interface{}) removeChannelEvent {
	return removeChannelEvent{
		channelID: events.Read(event, "id", ""),
		domainID:  events.Read(event, "domain", ""),
	}
}
