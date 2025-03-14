// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package events

type createClientEvent struct {
	id          string
	opcuaNodeID string
}

type removeClientEvent struct {
	id string
}

type connectClientEvent struct {
	chanID    string
	clientIDs []string
}

type createChannelEvent struct {
	id             string
	opcuaServerURI string
}

type removeChannelEvent struct {
	id string
}
