package serial

type EventType int

// TODO need to export this as a shared package and use it in both projects
const (
	Error    = iota
	Latency  = iota
	Ping     = iota
	Update   = iota
	Position = iota
)
