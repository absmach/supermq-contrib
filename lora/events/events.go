// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package events

type createClientEvent struct {
	id         string
	loraDevEUI string
}

type removeClientEvent struct {
	id string
}

type createChannelEvent struct {
	id        string
	loraAppID string
}

type removeChannelEvent struct {
	id string
}

type connectionClientEvent struct {
	chanID    string
	clientIDs []string
}
