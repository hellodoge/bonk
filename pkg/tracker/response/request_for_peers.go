package response

import (
	"github.com/hellodoge/bonk/bonk"
	"time"
)

type RequestForPeersResponse struct {
	Interval time.Duration
	Peers    []bonk.Peer
}
