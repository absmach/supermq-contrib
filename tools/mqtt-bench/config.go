// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package bench

// Keep struct names exported, otherwise Viper unmarshalling won't work.
type mqttBrokerConfig struct {
	URL string `toml:"url" mapstructure:"url"`
}

type mqttMessageConfig struct {
	Size    int    `toml:"size" mapstructure:"size"`
	Payload string `toml:"payload" mapstructure:"payload"`
	Format  string `toml:"format" mapstructure:"format"`
	QoS     int    `toml:"qos" mapstructure:"qos"`
	Retain  bool   `toml:"retain" mapstructure:"retain"`
}

type mqttTLSConfig struct {
	MTLS       bool   `toml:"mtls" mapstructure:"mtls"`
	SkipTLSVer bool   `toml:"skiptlsver" mapstructure:"skiptlsver"`
	CA         string `toml:"ca" mapstructure:"ca"`
}

type mqttConfig struct {
	Broker  mqttBrokerConfig  `toml:"broker" mapstructure:"broker"`
	Message mqttMessageConfig `toml:"message" mapstructure:"message"`
	Timeout int               `toml:"timeout" mapstructure:"timeout"`
	TLS     mqttTLSConfig     `toml:"tls" mapstructure:"tls"`
}

type testConfig struct {
	Count int `toml:"count" mapstructure:"count"`
	Pubs  int `toml:"pubs" mapstructure:"pubs"`
	Subs  int `toml:"subs" mapstructure:"subs"`
}

type logConfig struct {
	Quiet bool `toml:"quiet" mapstructure:"quiet"`
}

type supermqFile struct {
	ConnFile string `toml:"connections_file" mapstructure:"connections_file"`
}

type mgClient struct {
	ClientID  string `toml:"client_id" mapstructure:"client_id"`
	ClientKey string `toml:"client_key" mapstructure:"client_key"`
	MTLSCert  string `toml:"mtls_cert" mapstructure:"mtls_cert"`
	MTLSKey   string `toml:"mtls_key" mapstructure:"mtls_key"`
}

type mgChannel struct {
	ChannelID string `toml:"channel_id" mapstructure:"channel_id"`
}

type supermq struct {
	Clients  []mgClient  `toml:"clients" mapstructure:"clients"`
	Channels []mgChannel `toml:"channels" mapstructure:"channels"`
}

// Config struct holds benchmark configuration.
type Config struct {
	MQTT mqttConfig  `toml:"mqtt" mapstructure:"mqtt"`
	Test testConfig  `toml:"test" mapstructure:"test"`
	Log  logConfig   `toml:"log" mapstructure:"log"`
	Mg   supermqFile `toml:"supermq" mapstructure:"supermq"`
}
