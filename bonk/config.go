package bonk

import "time"

type Config struct {
	PeerID         PeerID
	Port           uint16
	TrackerTimeout time.Duration
	PeerTimeout    time.Duration
	RequestCompact bool
}
