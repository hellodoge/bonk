package torrent

import "github.com/hellodoge/bonk/bonk"

func StringToPeerID(str string) bonk.PeerID {
	var peerId bonk.PeerID
	copy(peerId[:], str)
	return peerId
}
