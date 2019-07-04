package model

// ClientConfig has client configuration configured by server.
type ClientConfig struct {
	PingInterval int    `json:"ping_interval"`
	DeviceID     string `json:"device_id"`
}
