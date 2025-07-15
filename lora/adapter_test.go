// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package lora_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/absmach/supermq-contrib/lora"
	"github.com/absmach/supermq-contrib/lora/mocks"
	"github.com/absmach/supermq/pkg/errors"
	pubmocks "github.com/absmach/supermq/pkg/messaging/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	clientID  = "clientID-1"
	chanID    = "chanID-1"
	devEUI    = "devEUI-1"
	appID     = "appID-1"
	clientID2 = "clientID-2"
	chanID2   = "chanID-2"
	devEUI2   = "devEUI-2"
	appID2    = "appID-2"
	msg       = `[{"bn":"msg-base-name","n":"temperature","v": 17},{"n":"humidity","v": 56}]`
	invalid   = "wrong"
	domainID  = "domainID-1"
)

var (
	pub                            *pubmocks.PubSub
	clientsRM, channelsRM, connsRM *mocks.RouteMapRepository
)

func newService() lora.Service {
	pub = new(pubmocks.PubSub)
	clientsRM = new(mocks.RouteMapRepository)
	channelsRM = new(mocks.RouteMapRepository)
	connsRM = new(mocks.RouteMapRepository)

	return lora.New(pub, clientsRM, channelsRM, connsRM)
}

func TestPublish(t *testing.T) {
	svc := newService()

	msgBase64 := base64.StdEncoding.EncodeToString([]byte(msg))

	cases := []struct {
		desc           string
		err            error
		msg            lora.Message
		getClientErr   error
		getChannelErr  error
		connectionsErr error
		publishErr     error
	}{
		{
			desc: "publish message with existing route-map and valid Data",
			err:  nil,
			msg: lora.Message{
				DeviceInfo: lora.DeviceInfo{
					ApplicationID: appID,
					DevEUI:        devEUI,
				},
				Data: msgBase64,
			},
			getClientErr:   nil,
			getChannelErr:  nil,
			connectionsErr: nil,
			publishErr:     nil,
		},
		{
			desc: "publish message with existing route-map and invalid Data",
			err:  lora.ErrMalformedMessage,
			msg: lora.Message{
				DeviceInfo: lora.DeviceInfo{
					ApplicationID: appID,
					DevEUI:        devEUI,
				},
				Data: "wrong",
			},
			getClientErr:   nil,
			getChannelErr:  nil,
			connectionsErr: nil,
			publishErr:     errors.New("Failed publishing"),
		},
		{
			desc: "publish message with non existing appID route-map",
			err:  lora.ErrNotFoundApp,
			msg: lora.Message{
				DeviceInfo: lora.DeviceInfo{
					ApplicationID: "wrong",
					DevEUI:        devEUI,
				},
			},
			getChannelErr: lora.ErrNotFoundApp,
		},
		{
			desc: "publish message with non existing devEUI route-map",
			err:  lora.ErrNotFoundDev,
			msg: lora.Message{
				DeviceInfo: lora.DeviceInfo{
					ApplicationID: appID,
					DevEUI:        "wrong",
				},
			},
			getClientErr: lora.ErrNotFoundDev,
		},
		{
			desc: "publish message with non existing connection route-map",
			err:  lora.ErrNotConnected,
			msg: lora.Message{
				DeviceInfo: lora.DeviceInfo{
					ApplicationID: appID2,
					DevEUI:        devEUI2,
				},
			},
			connectionsErr: lora.ErrNotConnected,
		},
		{
			desc: "publish message with wrong Object",
			err:  errors.New("json: unsupported type: chan int"),
			msg: lora.Message{
				DeviceInfo: lora.DeviceInfo{
					ApplicationID: appID2,
					DevEUI:        devEUI2,
				},
				Object: make(chan int),
			},
		},
		{
			desc: "publish message with valid Object",
			err:  nil,
			msg: lora.Message{
				DeviceInfo: lora.DeviceInfo{
					ApplicationID: appID2,
					DevEUI:        devEUI2,
				},
				Object: map[string]interface{}{"key": "value"},
			},
		},
	}

	for _, tc := range cases {
		repoCall := clientsRM.On("Get", context.Background(), tc.msg.DeviceInfo.DevEUI).Return(tc.msg.DeviceInfo.DevEUI, tc.getClientErr)
		repoCall1 := channelsRM.On("Get", context.Background(), tc.msg.DeviceInfo.ApplicationID).Return(tc.msg.DeviceInfo.ApplicationID, tc.getChannelErr)
		repoCall2 := connsRM.On("Get", context.Background(), mock.Anything).Return("", tc.connectionsErr)
		repoCall3 := pub.On("Publish", context.Background(), tc.msg.DeviceInfo.ApplicationID, mock.Anything).Return(tc.publishErr)
		err := svc.Publish(context.Background(), &tc.msg)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
		repoCall1.Unset()
		repoCall2.Unset()
		repoCall3.Unset()
	}
}

func TestCreateChannel(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc     string
		err      error
		chanID   string
		domainID string
		appID    string
	}{
		{
			desc:     "create channel with valid data",
			err:      nil,
			chanID:   chanID,
			domainID: domainID,
			appID:    appID,
		},
		{
			desc:     "create channel with empty chanID and domainID",
			err:      lora.ErrNotFoundApp,
			chanID:   "",
			domainID: "",
			appID:    appID,
		},
		{
			desc:     "create channel with empty appID",
			err:      lora.ErrNotFoundApp,
			chanID:   chanID,
			domainID: domainID,
			appID:    "",
		},
	}

	for _, tc := range cases {
		repoCall := channelsRM.On("Save", context.Background(), tc.chanID+":"+tc.domainID, tc.appID).Return(tc.err)
		err := svc.CreateChannel(context.Background(), tc.chanID, tc.domainID, tc.appID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestCreateClient(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc     string
		err      error
		ClientID string
		DevEUI   string
	}{
		{
			desc:     "create client with valid data",
			err:      nil,
			ClientID: clientID,
			DevEUI:   devEUI,
		},
		{
			desc:     "create client with empty clientID",
			err:      lora.ErrNotFoundDev,
			ClientID: "",
			DevEUI:   devEUI,
		},
		{
			desc:     "create client with empty devEUI",
			err:      lora.ErrNotFoundDev,
			ClientID: clientID,
			DevEUI:   "",
		},
	}

	for _, tc := range cases {
		repoCall := clientsRM.On("Save", context.Background(), tc.ClientID, tc.DevEUI).Return(tc.err)
		err := svc.CreateClient(context.Background(), tc.ClientID, tc.DevEUI)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestConnectClient(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc          string
		err           error
		channelID     string
		clientID      string
		domainID      string
		getClientErr  error
		getChannelErr error
	}{
		{
			desc:          "connect client with valid data",
			err:           nil,
			channelID:     chanID,
			clientID:      clientID,
			domainID:      domainID,
			getClientErr:  nil,
			getChannelErr: nil,
		},
		{
			desc:         "connect client with non existing client",
			err:          lora.ErrNotFoundDev,
			channelID:    chanID,
			clientID:     invalid,
			domainID:     domainID,
			getClientErr: lora.ErrNotFoundDev,
		},
		{
			desc:          "connect client with non existing channel",
			err:           lora.ErrNotFoundApp,
			channelID:     invalid,
			clientID:      clientID,
			domainID:      domainID,
			getChannelErr: lora.ErrNotFoundApp,
		},
	}

	for _, tc := range cases {
		repoCall := clientsRM.On("Get", context.Background(), tc.clientID).Return(devEUI, tc.getClientErr)
		repoCall1 := channelsRM.On("Get", context.Background(), tc.channelID+":"+tc.domainID).Return(appID, tc.getChannelErr)
		repoCall2 := connsRM.On("Save", context.Background(), mock.Anything, mock.Anything).Return(tc.err)
		err := svc.ConnectClient(context.Background(), tc.channelID, tc.domainID, tc.clientID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
		repoCall1.Unset()
		repoCall2.Unset()
	}
}

func TestDisconnectClient(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc          string
		err           error
		channelID     string
		clientID      string
		domainID      string
		getClientErr  error
		getChannelErr error
	}{
		{
			desc:          "disconnect client with valid data",
			err:           nil,
			channelID:     chanID,
			clientID:      clientID,
			domainID:      domainID,
			getClientErr:  nil,
			getChannelErr: nil,
		},
		{
			desc:         "disconnect client with non existing client ID",
			err:          lora.ErrNotFoundDev,
			channelID:    chanID,
			clientID:     invalid,
			domainID:     domainID,
			getClientErr: lora.ErrNotFoundDev,
		},
		{
			desc:          "disconnect client with non existing channel",
			err:           lora.ErrNotFoundApp,
			channelID:     invalid,
			clientID:      clientID,
			domainID:      domainID,
			getChannelErr: lora.ErrNotFoundApp,
		},
	}

	for _, tc := range cases {
		repoCall := clientsRM.On("Get", context.Background(), tc.clientID).Return(devEUI, tc.getClientErr)
		repoCall1 := channelsRM.On("Get", context.Background(), tc.channelID+":"+tc.domainID).Return(appID, tc.getChannelErr)
		repoCall2 := connsRM.On("Remove", context.Background(), mock.Anything).Return(tc.err)
		err := svc.DisconnectClient(context.Background(), tc.channelID, tc.domainID, tc.clientID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
		repoCall1.Unset()
		repoCall2.Unset()
	}
}

func TestRemoveChannel(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc     string
		err      error
		chanID   string
		domainID string
	}{
		{
			desc:     "remove channel with valid data",
			err:      nil,
			chanID:   chanID,
			domainID: domainID,
		},
		{
			desc:     "remove channel with non existing channel",
			err:      lora.ErrNotFoundApp,
			chanID:   invalid,
			domainID: domainID,
		},
		{
			desc:     "remove channel with empty channelID",
			err:      lora.ErrNotFoundApp,
			chanID:   "",
			domainID: domainID,
		},
		{
			desc:     "remove channel with empty domainID",
			err:      lora.ErrNotFoundApp,
			chanID:   chanID,
			domainID: "",
		},
	}

	for _, tc := range cases {
		repoCall := channelsRM.On("Remove", context.Background(), tc.chanID+":"+tc.domainID).Return(tc.err)
		err := svc.RemoveChannel(context.Background(), tc.chanID, tc.domainID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestRemoveClient(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc     string
		err      error
		ClientID string
	}{
		{
			desc:     "remove client with valid data",
			err:      nil,
			ClientID: clientID,
		},
		{
			desc:     "remove client with non existing client",
			err:      lora.ErrNotFoundDev,
			ClientID: invalid,
		},
		{
			desc:     "remove client with empty clientID",
			err:      lora.ErrNotFoundDev,
			ClientID: "",
		},
	}

	for _, tc := range cases {
		repoCall := clientsRM.On("Remove", context.Background(), tc.ClientID).Return(tc.err)
		err := svc.RemoveClient(context.Background(), tc.ClientID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestUpdateChannel(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc     string
		err      error
		chanID   string
		domainID string
		appID    string
	}{
		{
			desc:     "update channel with valid data",
			err:      nil,
			chanID:   chanID,
			domainID: domainID,
			appID:    appID,
		},
		{
			desc:     "update channel with non existing channel",
			err:      lora.ErrNotFoundApp,
			chanID:   invalid,
			domainID: domainID,
			appID:    appID,
		},
		{
			desc:     "update channel with empty channelID",
			err:      lora.ErrNotFoundApp,
			chanID:   "",
			domainID: domainID,
			appID:    appID,
		},
		{
			desc:     "update channel with empty appID",
			err:      lora.ErrNotFoundApp,
			chanID:   chanID,
			domainID: domainID,
			appID:    "",
		},
		{
			desc:     "update channel with non existing appID",
			err:      lora.ErrNotFoundApp,
			chanID:   chanID,
			domainID: domainID,
			appID:    invalid,
		},
	}

	for _, tc := range cases {
		repoCall := channelsRM.On("Save", context.Background(), tc.chanID+":"+tc.domainID, tc.appID).Return(tc.err)
		err := svc.UpdateChannel(context.Background(), tc.chanID, tc.domainID, tc.appID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}

func TestUpdateClient(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc     string
		err      error
		ClientID string
		DevEUI   string
	}{
		{
			desc:     "update client with valid data",
			err:      nil,
			ClientID: clientID,
			DevEUI:   devEUI,
		},
		{
			desc:     "update client with non existing client",
			err:      lora.ErrNotFoundDev,
			ClientID: invalid,
			DevEUI:   devEUI,
		},
		{
			desc:     "update client with empty clientID",
			err:      lora.ErrNotFoundDev,
			ClientID: "",
			DevEUI:   devEUI,
		},
		{
			desc:     "update client with empty devEUI",
			err:      lora.ErrNotFoundDev,
			ClientID: clientID,
			DevEUI:   "",
		},
		{
			desc:     "update client with non existing devEUI",
			err:      lora.ErrNotFoundDev,
			ClientID: clientID,
			DevEUI:   invalid,
		},
	}

	for _, tc := range cases {
		repoCall := clientsRM.On("Save", context.Background(), tc.ClientID, tc.DevEUI).Return(tc.err)
		err := svc.UpdateClient(context.Background(), tc.ClientID, tc.DevEUI)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		repoCall.Unset()
	}
}
