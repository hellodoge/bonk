package parser

import (
	"errors"
	"github.com/hellodoge/bonk/pkg/errors/wrapped"
)

var (
	ErrInvalidPiecesHashesLength = errors.New("invalid length of pieces hashes")
	ErrParsingError              = errors.New("error occurred while parsing")
)

func newParsingError(inner error) wrapped.Error {
	return wrapped.Error{
		Outer: ErrParsingError,
		Inner: inner,
	}
}
