// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package events

import (
	"context"
	"errors"

	"github.com/absmach/supermq-contrib/opcua"
	"github.com/absmach/supermq/pkg/events"
)

const (
	keyType      = "opcua"
	keyNodeID    = "node_id"
	keyServerURI = "server_uri"

	clientPrefix = "client."
	clientCreate = clientPrefix + "create"
	clientUpdate = clientPrefix + "update"
	clientRemove = clientPrefix + "remove"

	channelPrefix     = "channel."
	channelCreate     = channelPrefix + "create"
	channelUpdate     = channelPrefix + "update"
	channelRemove     = channelPrefix + "remove"
	channelConnect    = channelPrefix + "assign"
	channelDisconnect = channelPrefix + "unassign"
)

var (
	errMetadataType = errors.New("metadatada is not of type opcua")

	errMetadataFormat = errors.New("malformed metadata")

	errMetadataServerURI = errors.New("ServerURI not found in channel metadatada")

	errMetadataNodeID = errors.New("NodeID not found in client metadatada")
)

type eventHandler struct {
	svc opcua.Service
}

// NewEventHandler returns new event store handler.
func NewEventHandler(svc opcua.Service) events.EventHandler {
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
	case clientCreate:
		cte, e := decodeCreateClient(msg)
		if e != nil {
			err = e
			break
		}
		err = es.svc.CreateClient(ctx, cte.id, cte.opcuaNodeID)
	case clientUpdate:
		ute, e := decodeCreateClient(msg)
		if e != nil {
			err = e
			break
		}
		err = es.svc.CreateClient(ctx, ute.id, ute.opcuaNodeID)
	case clientRemove:
		rte := decodeRemoveClient(msg)
		err = es.svc.RemoveClient(ctx, rte.id)
	case channelCreate:
		cce, e := decodeCreateChannel(msg)
		if e != nil {
			err = e
			break
		}
		err = es.svc.CreateChannel(ctx, cce.id, cce.opcuaServerURI)
	case channelUpdate:
		uce, e := decodeCreateChannel(msg)
		if e != nil {
			err = e
			break
		}
		err = es.svc.CreateChannel(ctx, uce.id, uce.opcuaServerURI)
	case channelRemove:
		rce := decodeRemoveChannel(msg)
		err = es.svc.RemoveChannel(ctx, rce.id)
	case channelConnect:
		rce := decodeConnectClient(msg)
		err = es.svc.ConnectClient(ctx, rce.chanID, rce.clientIDs)
	case channelDisconnect:
		rce := decodeDisconnectClient(msg)
		err = es.svc.DisconnectClient(ctx, rce.chanID, rce.clientIDs)
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

	metadataOpcua, ok := metadata[keyType]
	if !ok {
		return createClientEvent{}, errMetadataType
	}

	metadataVal, ok := metadataOpcua.(map[string]interface{})
	if !ok {
		return createClientEvent{}, errMetadataFormat
	}

	val, ok := metadataVal[keyNodeID].(string)
	if !ok || val == "" {
		return createClientEvent{}, errMetadataNodeID
	}

	cte.opcuaNodeID = val
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

	metadataOpcua, ok := metadata[keyType]
	if !ok {
		return createChannelEvent{}, errMetadataType
	}

	metadataVal, ok := metadataOpcua.(map[string]interface{})
	if !ok {
		return createChannelEvent{}, errMetadataFormat
	}

	val, ok := metadataVal[keyServerURI].(string)
	if !ok || val == "" {
		return createChannelEvent{}, errMetadataServerURI
	}

	cce.opcuaServerURI = val
	return cce, nil
}

func decodeRemoveChannel(event map[string]interface{}) removeChannelEvent {
	return removeChannelEvent{
		id: events.Read(event, "id", ""),
	}
}

func decodeConnectClient(event map[string]interface{}) connectClientEvent {
	return connectClientEvent{
		chanID:    events.Read(event, "group_id", ""),
		clientIDs: events.ReadStringSlice(event, "member_ids"),
	}
}

func decodeDisconnectClient(event map[string]interface{}) connectClientEvent {
	return connectClientEvent{
		chanID:    events.Read(event, "group_id", ""),
		clientIDs: events.ReadStringSlice(event, "member_ids"),
	}
}
