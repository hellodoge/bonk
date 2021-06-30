package peer

import (
	"errors"
	"github.com/hellodoge/bonk/internal/errors/wrapped"
)

var (
	ErrConnectingToPeer = errors.New("error connecting to peer")
	ErrWritingToPeer    = errors.New("error writing to peer")
	ErrReadingFromPeer  = errors.New("error reading from peer")
	ErrGettingHandshake = errors.New("error getting handshake")
)

func newConnectingToPeerError(inner error) wrapped.Error {
	return wrapped.Error{
		Outer: ErrConnectingToPeer,
		Inner: inner,
	}
}

func newWritingToPeerError(inner error) wrapped.Error {
	return wrapped.Error{
		Outer: ErrWritingToPeer,
		Inner: inner,
	}
}

func newReadingFromPeerError(inner error) wrapped.Error {
	return wrapped.Error{
		Outer: ErrReadingFromPeer,
		Inner: inner,
	}
}

func newGettingHandshakeError(inner error) wrapped.Error {
	return wrapped.Error{
		Outer: ErrGettingHandshake,
		Inner: inner,
	}
}
