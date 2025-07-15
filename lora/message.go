// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package lora

// RxInfo receiver parameters.
type RxInfo []struct {
	GatewayID string            `json:"gatewayId"`
	UplinkID  uint64            `json:"uplinkId"`
	Context   string            `json:"context"`
	Rssi      float64           `json:"rssi"`
	Snr       float64           `json:"snr"`
	Metadata  map[string]string `json:"metadata"`
}

// TxInfo transmeter parameters.
type TxInfo struct {
	Frequency  float64                `json:"frequency"`
	Modulation map[string]interface{} `json:"modulation"`
}

type DeviceInfo struct {
	TenantID           string            `json:"tenantId"`
	TenantName         string            `json:"tenantName"`
	ApplicationID      string            `json:"applicationId"`
	ApplicationName    string            `json:"applicationName"`
	DeviceProfileID    string            `json:"deviceProfileId"`
	DeviceProfileName  string            `json:"deviceProfileName"`
	DeviceName         string            `json:"deviceName"`
	DevEUI             string            `json:"devEui"`
	DeviceClassEnabled string            `json:"deviceClassEnabled"`
	Tags               map[string]string `json:"tags"`
}

// Message lora msg (https://www.chirpstack.io/docs/chirpstack/integrations/events.html).
type Message struct {
	DeduplicationID string     `json:"deduplicationId"`
	Time            string     `json:"time"`
	DeviceInfo      DeviceInfo `json:"deviceInfo"`
	DevAddr         string     `json:"devAddr"`
	Adr             bool       `json:"adr"`
	Dr              int        `json:"dr"`
	FCnt            int        `json:"fCnt"`
	FPort           int        `json:"fPort"`
	Confirmed       bool       `json:"confirmed"`
	Data            string     `json:"data"`
	RxInfo          RxInfo     `json:"rxInfo"`
	TxInfo          TxInfo     `json:"txInfo"`
	RegionConfigID  string     `json:"regionConfigId"`
	Object          any        `json:"object"`
}
