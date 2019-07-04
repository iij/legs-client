package model

import "github.com/iij/legs-client/daemon/model/status"

// Status has daemon and WebSocket connection status
type Status struct {
	Daemon   status.Daemon     `json:"daemon_status"`
	Conn     status.Connection `json:"connection_status"`
	DeviceID string            `json:"device_id"`
}
