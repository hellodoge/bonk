package bonk

import "net"

const PeerIdLength = 20

type PeerID [PeerIdLength]byte

type Peer struct {
	ID   PeerID
	Addr net.TCPAddr
}
