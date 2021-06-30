package bonk

import "crypto/sha1"

type Piece struct {
	Index   uint32
	Content []byte
}

type PieceInfo struct {
	Index  uint32
	Length uint
	Hash   [sha1.Size]byte
}

type Part struct {
	Index  uint32
	Offset uint32
	Block  []byte
}

type Request struct {
	Index     uint32
	Begin     uint32
	Length    uint32
	Cancelled bool
}
