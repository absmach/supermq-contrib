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
	channelID string
	loraAppID string
	domainID  string
}

type removeChannelEvent struct {
	channelID string
	domainID  string
}

type connectionEvent struct {
	chanIDs   []string
	clientIDs []string
	domainID  string
}
