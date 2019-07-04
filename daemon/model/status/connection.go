package status

// Connection is WebSocket connection status
type Connection int

// Disconnected mean device is disconnected
// Connected mean device is connected
const (
	Disconnected Connection = iota
	Connected
)

func (c Connection) String() string {
	switch c {
	case Disconnected:
		return "Disconnected"
	case Connected:
		return "Connected"
	default:
		return "Unknown"
	}
}
