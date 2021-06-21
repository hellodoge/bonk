package delivery

import (
	"errors"
	"github.com/hellodoge/bonk/internal/errors/wrapped"
)

var (
	ErrParsingTrackerURL = errors.New("error occurred while parsing tracker url")
	ErrMakingRequest     = errors.New("error occurred while making request to tracker")
	ErrTimeout           = errors.New("timeout error while making request to tracker")
)

func newParsingTrackerURLError(inner error) wrapped.Error {
	return wrapped.Error{
		Outer: ErrParsingTrackerURL,
		Inner: inner,
	}
}

func newMakingRequestError(inner error) wrapped.Error {
	return wrapped.Error{
		Outer: ErrMakingRequest,
		Inner: inner,
	}
}

func newTimeoutError(inner error) wrapped.Error {
	return wrapped.Error{
		Outer: ErrTimeout,
		Inner: inner,
	}
}
