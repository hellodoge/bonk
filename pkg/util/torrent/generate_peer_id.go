package torrent

import (
	"encoding/base64"
	"github.com/hellodoge/bonk/bonk"
	"math/rand" // no need of cryptographic security
)

func GenerateRandomPeerID() ([bonk.PeerIdLength]byte, error) {
	const NumOfBytesToGenerate = bonk.PeerIdLength * 6 /*bits per base64 char*/ / 8 /*bits in byte*/
	var buffer = make([]byte, NumOfBytesToGenerate)
	_, err := rand.Read(buffer)
	if err != nil {
		return [bonk.PeerIdLength]byte{}, err
	}
	var result [bonk.PeerIdLength]byte
	base64.URLEncoding.Encode(result[:], buffer)
	return result, nil
}
