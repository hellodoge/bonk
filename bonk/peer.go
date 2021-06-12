package bonk

import "net"

const PeerIdLength = 20

type Peer struct {
	ID string
	Addr net.TCPAddr
}