package messages

import (
	"github.com/hellodoge/bonk/bonk"
	"github.com/hellodoge/bonk/pkg/bitfield"
)

type Choke struct{}
type Unchoke struct{}
type Interested struct{}
type NotInterested struct{}

type Have struct {
	Index uint32
}

type BitField struct {
	Bitfield bitfield.BitField
}

type Request struct {
	Request bonk.Request
}

type Piece struct {
	Piece bonk.Part
}

type Cancel struct {
	Request
}
