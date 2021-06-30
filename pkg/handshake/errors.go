package handshake

import (
	"errors"
	wrapped2 "github.com/hellodoge/bonk/internal/errors/wrapped"
)

var (
	ErrReceivingHandshake = errors.New("error occurred while receiving handshake")
	ErrParsingHandshake   = errors.New("parsing handshake error")
)

func newReceivingHandshakeError(inner error) wrapped2.Error {
	return wrapped2.Error{
		Outer: ErrReceivingHandshake,
		Inner: inner,
	}
}

func newParsingHandshakeError(inner error) wrapped2.Error {
	return wrapped2.Error{
		Outer: ErrParsingHandshake,
		Inner: inner,
	}
}
