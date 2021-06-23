package messages

import (
	"errors"
	wrapped2 "github.com/hellodoge/bonk/internal/errors/wrapped"
)

var (
	ErrReceivingMessage = errors.New("error occurred while receiving message")
	ErrParsingMessage   = errors.New("error occurred while parsing handshake")
)

func newReceivingMessageError(inner error) wrapped2.Error {
	return wrapped2.Error{
		Outer: ErrReceivingMessage,
		Inner: inner,
	}
}

func newParsingMessageError(inner error) wrapped2.Error {
	return wrapped2.Error{
		Outer: ErrParsingMessage,
		Inner: inner,
	}
}
