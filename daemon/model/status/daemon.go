package status

// Daemon is daemon status
type Daemon int

// Stopped mean daemon is stopped
// Startted mean daemon is startted
const (
	Stopped Daemon = iota
	Startted
)

func (d Daemon) String() string {
	switch d {
	case Stopped:
		return "Stopped"
	case Startted:
		return "Startted"
	default:
		return "Unknown"
	}
}
