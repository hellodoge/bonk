package parser

import (
	"errors"
	"github.com/hellodoge/bonk/pkg/errors/wrapped"
)

var (
	ErrParsingResponse = errors.New("error while parsing response from tracker")
)

func newParsingResponseError(inner error) wrapped.Error {
	return wrapped.Error{
		Outer: ErrParsingResponse,
		Inner: inner,
	}
}